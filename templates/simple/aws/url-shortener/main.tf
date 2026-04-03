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

# --- DynamoDB table for short URLs ---

resource "aws_dynamodb_table" "urls" {
  name         = "${var.project_name}-urls"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "short_code"

  attribute {
    name = "short_code"
    type = "S"
  }

  ttl {
    attribute_name = "expires_at"
    enabled        = true
  }

  server_side_encryption {
    enabled = true
  }

  tags = {
    Name      = "${var.project_name}-urls"
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

# --- Lambda ---

resource "aws_iam_role" "lambda" {
  name = "${var.project_name}-shortener-lambda-role"

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
    Name      = "${var.project_name}-shortener-lambda-role"
    ManagedBy = "aeroform"
  }
}

resource "aws_iam_role_policy" "lambda" {
  name = "${var.project_name}-shortener-lambda-policy"
  role = aws_iam_role.lambda.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["dynamodb:GetItem", "dynamodb:PutItem"]
        Resource = aws_dynamodb_table.urls.arn
      },
      {
        Effect   = "Allow"
        Action   = ["logs:CreateLogGroup", "logs:CreateLogStream", "logs:PutLogEvents"]
        Resource = "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:*"
      }
    ]
  })
}

resource "aws_lambda_function" "shortener" {
  function_name = "${var.project_name}-shortener"
  role          = aws_iam_role.lambda.arn
  handler       = "index.handler"
  runtime       = "nodejs20.x"
  timeout       = 10
  memory_size   = 128

  filename         = data.archive_file.lambda_placeholder.output_path
  source_code_hash = data.archive_file.lambda_placeholder.output_base64sha256

  environment {
    variables = {
      TABLE_NAME = aws_dynamodb_table.urls.name
    }
  }

  tags = {
    Name      = "${var.project_name}-shortener"
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

data "archive_file" "lambda_placeholder" {
  type        = "zip"
  output_path = "${path.module}/lambda_placeholder.zip"

  source {
    content  = <<-JS
      const { DynamoDBClient } = require('@aws-sdk/client-dynamodb');
      const { DynamoDBDocumentClient, PutCommand, GetCommand } = require('@aws-sdk/lib-dynamodb');
      const client = DynamoDBDocumentClient.from(new DynamoDBClient({}));
      const TABLE = process.env.TABLE_NAME;
      const headers = { 'Content-Type': 'application/json', 'Access-Control-Allow-Origin': '*' };

      exports.handler = async (event) => {
        const method = event.requestContext?.http?.method || event.httpMethod;
        try {
          if (method === 'POST') {
            const body = JSON.parse(event.body || '{}');
            if (!body.url) return { statusCode: 400, headers, body: JSON.stringify({ error: 'url is required' }) };
            const code = Date.now().toString(36) + Math.random().toString(36).slice(2, 5);
            await client.send(new PutCommand({ TableName: TABLE, Item: { short_code: code, long_url: body.url, created_at: Date.now() } }));
            return { statusCode: 201, headers, body: JSON.stringify({ short_code: code, url: body.url }) };
          }
          if (method === 'GET') {
            const code = event.pathParameters?.code;
            if (!code) return { statusCode: 400, headers, body: JSON.stringify({ error: 'code is required' }) };
            const { Item } = await client.send(new GetCommand({ TableName: TABLE, Key: { short_code: code } }));
            if (!Item) return { statusCode: 404, headers, body: JSON.stringify({ error: 'not found' }) };
            return { statusCode: 301, headers: { ...headers, Location: Item.long_url }, body: '' };
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

resource "aws_cloudwatch_log_group" "shortener" {
  name              = "/aws/lambda/${aws_lambda_function.shortener.function_name}"
  retention_in_days = 7

  tags = { ManagedBy = "aeroform" }
}

# --- API Gateway ---

resource "aws_apigatewayv2_api" "shortener" {
  name          = "${var.project_name}-shortener-api"
  protocol_type = "HTTP"

  cors_configuration {
    allow_origins = ["*"]
    allow_methods = ["GET", "POST", "OPTIONS"]
    allow_headers = ["Content-Type"]
    max_age       = 300
  }

  tags = {
    Name      = "${var.project_name}-shortener-api"
    ManagedBy = "aeroform"
  }
}

resource "aws_apigatewayv2_integration" "shortener" {
  api_id                 = aws_apigatewayv2_api.shortener.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.shortener.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "shorten" {
  api_id    = aws_apigatewayv2_api.shortener.id
  route_key = "POST /shorten"
  target    = "integrations/${aws_apigatewayv2_integration.shortener.id}"
}

resource "aws_apigatewayv2_route" "redirect" {
  api_id    = aws_apigatewayv2_api.shortener.id
  route_key = "GET /{code}"
  target    = "integrations/${aws_apigatewayv2_integration.shortener.id}"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.shortener.id
  name        = "$default"
  auto_deploy = true

  tags = { ManagedBy = "aeroform" }
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.shortener.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.shortener.execution_arn}/*/*"
}
