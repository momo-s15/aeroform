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

# --- S3 bucket for static content ---

resource "aws_s3_bucket" "site" {
  bucket        = "${var.project_name}-site"
  force_destroy = true

  tags = {
    Name      = var.project_name
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

resource "aws_s3_bucket_website_configuration" "site" {
  bucket = aws_s3_bucket.site.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }
}

resource "aws_s3_bucket_public_access_block" "site" {
  bucket = aws_s3_bucket.site.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_server_side_encryption_configuration" "site" {
  bucket = aws_s3_bucket.site.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_versioning" "site" {
  bucket = aws_s3_bucket.site.id

  versioning_configuration {
    status = "Enabled"
  }
}

# --- CloudFront OAI for secure S3 access ---

resource "aws_cloudfront_origin_access_identity" "site" {
  comment = "OAI for ${var.project_name}"
}

resource "aws_s3_bucket_policy" "site" {
  bucket = aws_s3_bucket.site.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "AllowCloudFrontOAI"
        Effect    = "Allow"
        Principal = { AWS = aws_cloudfront_origin_access_identity.site.iam_arn }
        Action    = "s3:GetObject"
        Resource  = "${aws_s3_bucket.site.arn}/*"
      }
    ]
  })

  depends_on = [aws_s3_bucket_public_access_block.site]
}

# --- CloudFront distribution ---

resource "aws_cloudfront_distribution" "site" {
  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"
  comment             = "${var.project_name} — Aeroform contact-form"
  price_class         = "PriceClass_100"
  wait_for_deployment = true

  origin {
    domain_name = aws_s3_bucket.site.bucket_regional_domain_name
    origin_id   = "s3-${var.project_name}"

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.site.cloudfront_access_identity_path
    }
  }

  default_cache_behavior {
    allowed_methods        = ["GET", "HEAD"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = "s3-${var.project_name}"
    viewer_protocol_policy = "redirect-to-https"
    compress               = true

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }

    min_ttl     = 0
    default_ttl = 86400
    max_ttl     = 31536000
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }

  tags = {
    Name      = var.project_name
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

# --- Lambda function for the contact form backend ---

resource "aws_iam_role" "lambda" {
  name = "${var.project_name}-lambda-role"

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
    Name      = "${var.project_name}-lambda-role"
    ManagedBy = "aeroform"
  }
}

resource "aws_iam_role_policy" "lambda_ses" {
  name = "${var.project_name}-lambda-ses"
  role = aws_iam_role.lambda.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["ses:SendEmail", "ses:SendRawEmail"]
        Resource = "*"
      },
      {
        Effect   = "Allow"
        Action   = ["logs:CreateLogGroup", "logs:CreateLogStream", "logs:PutLogEvents"]
        Resource = "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:*"
      }
    ]
  })
}

resource "aws_lambda_function" "contact" {
  function_name = "${var.project_name}-contact"
  role          = aws_iam_role.lambda.arn
  handler       = "index.handler"
  runtime       = "nodejs20.x"
  timeout       = 10
  memory_size   = 128

  filename         = data.archive_file.lambda_placeholder.output_path
  source_code_hash = data.archive_file.lambda_placeholder.output_base64sha256

  environment {
    variables = {
      RECIPIENT_EMAIL = var.contact_email
      PROJECT_NAME    = var.project_name
    }
  }

  tags = {
    Name      = "${var.project_name}-contact"
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

data "archive_file" "lambda_placeholder" {
  type        = "zip"
  output_path = "${path.module}/lambda_placeholder.zip"

  source {
    content  = <<-JS
      exports.handler = async (event) => {
        const body = JSON.parse(event.body || '{}');
        const AWS = require('@aws-sdk/client-ses');
        const ses = new AWS.SES();
        await ses.sendEmail({
          Source: process.env.RECIPIENT_EMAIL,
          Destination: { ToAddresses: [process.env.RECIPIENT_EMAIL] },
          Message: {
            Subject: { Data: 'Contact from ' + (body.name || 'visitor') + ' via ' + process.env.PROJECT_NAME },
            Body: { Text: { Data: 'Name: ' + (body.name || '') + '\nEmail: ' + (body.email || '') + '\nMessage: ' + (body.message || '') } }
          }
        });
        return { statusCode: 200, headers: { 'Access-Control-Allow-Origin': '*' }, body: JSON.stringify({ ok: true }) };
      };
    JS
    filename = "index.js"
  }
}

resource "aws_cloudwatch_log_group" "contact" {
  name              = "/aws/lambda/${aws_lambda_function.contact.function_name}"
  retention_in_days = 7

  tags = {
    ManagedBy = "aeroform"
  }
}

# --- API Gateway (HTTP API) ---

resource "aws_apigatewayv2_api" "contact" {
  name          = "${var.project_name}-contact-api"
  protocol_type = "HTTP"

  cors_configuration {
    allow_origins = ["*"]
    allow_methods = ["POST", "OPTIONS"]
    allow_headers = ["Content-Type"]
    max_age       = 300
  }

  tags = {
    Name      = "${var.project_name}-contact-api"
    ManagedBy = "aeroform"
  }
}

resource "aws_apigatewayv2_integration" "contact" {
  api_id                 = aws_apigatewayv2_api.contact.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.contact.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "contact" {
  api_id    = aws_apigatewayv2_api.contact.id
  route_key = "POST /contact"
  target    = "integrations/${aws_apigatewayv2_integration.contact.id}"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.contact.id
  name        = "$default"
  auto_deploy = true

  tags = {
    ManagedBy = "aeroform"
  }
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.contact.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.contact.execution_arn}/*/*"
}

# --- SES email identity ---

resource "aws_ses_email_identity" "contact" {
  email = var.contact_email
}
