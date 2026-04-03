terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

provider "aws" {
  region = var.region
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

# --- S3 bucket for uploaded files ---

resource "aws_s3_bucket" "uploads" {
  bucket        = "${var.project_name}-uploads"
  force_destroy = true

  tags = {
    Name      = "${var.project_name}-uploads"
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

resource "aws_s3_bucket_public_access_block" "uploads" {
  bucket = aws_s3_bucket.uploads.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_server_side_encryption_configuration" "uploads" {
  bucket = aws_s3_bucket.uploads.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_versioning" "uploads" {
  bucket = aws_s3_bucket.uploads.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_cors_configuration" "uploads" {
  bucket = aws_s3_bucket.uploads.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "GET"]
    allowed_origins = ["*"]
    max_age_seconds = 3600
  }
}

# --- Lambda for generating presigned URLs ---

resource "aws_iam_role" "lambda" {
  name = "${var.project_name}-upload-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = { Service = "lambda.amazonaws.com" }
        Action    = "sts:AssumeRole"
      }
    ]
  })

  tags = {
    Name      = "${var.project_name}-upload-lambda-role"
    ManagedBy = "aeroform"
  }
}

resource "aws_iam_role_policy" "lambda" {
  name = "${var.project_name}-upload-lambda-policy"
  role = aws_iam_role.lambda.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["s3:PutObject", "s3:GetObject"]
        Resource = "${aws_s3_bucket.uploads.arn}/*"
      },
      {
        Effect   = "Allow"
        Action   = ["logs:CreateLogGroup", "logs:CreateLogStream", "logs:PutLogEvents"]
        Resource = "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:*"
      }
    ]
  })
}

resource "aws_lambda_function" "upload" {
  function_name = "${var.project_name}-upload"
  role          = aws_iam_role.lambda.arn
  handler       = "index.handler"
  runtime       = "nodejs20.x"
  timeout       = 10
  memory_size   = 128

  filename         = data.archive_file.lambda_placeholder.output_path
  source_code_hash = data.archive_file.lambda_placeholder.output_base64sha256

  environment {
    variables = {
      BUCKET_NAME = aws_s3_bucket.uploads.id
    }
  }

  tags = {
    Name      = "${var.project_name}-upload"
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

data "archive_file" "lambda_placeholder" {
  type        = "zip"
  output_path = "${path.module}/lambda_placeholder.zip"

  source {
    content  = <<-JS
      const { S3Client, PutObjectCommand, GetObjectCommand } = require('@aws-sdk/client-s3');
      const { getSignedUrl } = require('@aws-sdk/s3-request-presigner');
      const s3 = new S3Client({});
      const BUCKET = process.env.BUCKET_NAME;
      const headers = { 'Content-Type': 'application/json', 'Access-Control-Allow-Origin': '*' };

      exports.handler = async (event) => {
        const method = event.requestContext?.http?.method || event.httpMethod;
        const body = JSON.parse(event.body || '{}');
        try {
          if (method === 'POST') {
            const key = body.filename || (Date.now().toString(36) + Math.random().toString(36).slice(2, 7));
            const url = await getSignedUrl(s3, new PutObjectCommand({ Bucket: BUCKET, Key: key, ContentType: body.contentType || 'application/octet-stream' }), { expiresIn: 3600 });
            return { statusCode: 200, headers, body: JSON.stringify({ uploadUrl: url, key }) };
          }
          if (method === 'GET') {
            const key = event.queryStringParameters?.key;
            if (!key) return { statusCode: 400, headers, body: JSON.stringify({ error: 'key is required' }) };
            const url = await getSignedUrl(s3, new GetObjectCommand({ Bucket: BUCKET, Key: key }), { expiresIn: 3600 });
            return { statusCode: 200, headers, body: JSON.stringify({ downloadUrl: url, key }) };
          }
          return { statusCode: 405, headers, body: JSON.stringify({ error: 'method not allowed' }) };
        } catch (err) {
          return { statusCode: 500, headers, body: JSON.stringify({ error: err.message }) };
        }
      };
    JS
    filename = "index.js"
  }
}

resource "aws_cloudwatch_log_group" "upload" {
  name              = "/aws/lambda/${aws_lambda_function.upload.function_name}"
  retention_in_days = 7

  tags = { ManagedBy = "aeroform" }
}

# --- API Gateway ---

resource "aws_apigatewayv2_api" "upload" {
  name          = "${var.project_name}-upload-api"
  protocol_type = "HTTP"

  cors_configuration {
    allow_origins = ["*"]
    allow_methods = ["GET", "POST", "OPTIONS"]
    allow_headers = ["Content-Type"]
    max_age       = 300
  }

  tags = {
    Name      = "${var.project_name}-upload-api"
    ManagedBy = "aeroform"
  }
}

resource "aws_apigatewayv2_integration" "upload" {
  api_id                 = aws_apigatewayv2_api.upload.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.upload.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "upload_post" {
  api_id    = aws_apigatewayv2_api.upload.id
  route_key = "POST /upload"
  target    = "integrations/${aws_apigatewayv2_integration.upload.id}"
}

resource "aws_apigatewayv2_route" "upload_get" {
  api_id    = aws_apigatewayv2_api.upload.id
  route_key = "GET /download"
  target    = "integrations/${aws_apigatewayv2_integration.upload.id}"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.upload.id
  name        = "$default"
  auto_deploy = true

  tags = { ManagedBy = "aeroform" }
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.upload.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.upload.execution_arn}/*/*"
}
