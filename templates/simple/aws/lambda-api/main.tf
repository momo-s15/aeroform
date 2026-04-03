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

# --- DynamoDB table (on-demand, free tier) ---

resource "aws_dynamodb_table" "main" {
  name         = "${var.project_name}-table"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }

  point_in_time_recovery {
    enabled = true
  }

  server_side_encryption {
    enabled = true
  }

  tags = {
    Name      = "${var.project_name}-table"
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

# --- Lambda execution role ---

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

resource "aws_iam_role_policy" "lambda" {
  name = "${var.project_name}-lambda-policy"
  role = aws_iam_role.lambda.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:Scan"
        ]
        Resource = aws_dynamodb_table.main.arn
      },
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:*"
      }
    ]
  })
}

# --- Lambda function ---

resource "aws_lambda_function" "api" {
  function_name = "${var.project_name}-api"
  role          = aws_iam_role.lambda.arn
  handler       = "index.handler"
  runtime       = "nodejs20.x"
  timeout       = 10
  memory_size   = 128

  filename         = data.archive_file.lambda_placeholder.output_path
  source_code_hash = data.archive_file.lambda_placeholder.output_base64sha256

  environment {
    variables = {
      TABLE_NAME   = aws_dynamodb_table.main.name
      PROJECT_NAME = var.project_name
    }
  }

  tags = {
    Name      = "${var.project_name}-api"
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
      const { DynamoDBDocumentClient, PutCommand, GetCommand, ScanCommand, DeleteCommand } = require('@aws-sdk/lib-dynamodb');
      const client = DynamoDBDocumentClient.from(new DynamoDBClient({}));
      const TABLE = process.env.TABLE_NAME;
      const headers = { 'Content-Type': 'application/json', 'Access-Control-Allow-Origin': '*' };

      exports.handler = async (event) => {
        const method = event.requestContext?.http?.method || event.httpMethod;
        const id = event.pathParameters?.id;
        try {
          if (method === 'GET' && !id) {
            const { Items } = await client.send(new ScanCommand({ TableName: TABLE }));
            return { statusCode: 200, headers, body: JSON.stringify(Items) };
          }
          if (method === 'GET' && id) {
            const { Item } = await client.send(new GetCommand({ TableName: TABLE, Key: { id } }));
            if (!Item) return { statusCode: 404, headers, body: JSON.stringify({ error: 'not found' }) };
            return { statusCode: 200, headers, body: JSON.stringify(Item) };
          }
          if (method === 'POST') {
            const body = JSON.parse(event.body || '{}');
            body.id = body.id || Date.now().toString(36) + Math.random().toString(36).slice(2, 7);
            await client.send(new PutCommand({ TableName: TABLE, Item: body }));
            return { statusCode: 201, headers, body: JSON.stringify(body) };
          }
          if (method === 'DELETE' && id) {
            await client.send(new DeleteCommand({ TableName: TABLE, Key: { id } }));
            return { statusCode: 200, headers, body: JSON.stringify({ deleted: id }) };
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

resource "aws_cloudwatch_log_group" "api" {
  name              = "/aws/lambda/${aws_lambda_function.api.function_name}"
  retention_in_days = 7

  tags = {
    ManagedBy = "aeroform"
  }
}

# --- API Gateway (HTTP API) ---

resource "aws_apigatewayv2_api" "api" {
  name          = "${var.project_name}-api"
  protocol_type = "HTTP"

  cors_configuration {
    allow_origins = ["*"]
    allow_methods = ["GET", "POST", "DELETE", "OPTIONS"]
    allow_headers = ["Content-Type"]
    max_age       = 300
  }

  tags = {
    Name      = "${var.project_name}-api"
    ManagedBy = "aeroform"
  }
}

resource "aws_apigatewayv2_integration" "api" {
  api_id                 = aws_apigatewayv2_api.api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.api.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "get_all" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "GET /items"
  target    = "integrations/${aws_apigatewayv2_integration.api.id}"
}

resource "aws_apigatewayv2_route" "get_one" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "GET /items/{id}"
  target    = "integrations/${aws_apigatewayv2_integration.api.id}"
}

resource "aws_apigatewayv2_route" "post" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "POST /items"
  target    = "integrations/${aws_apigatewayv2_integration.api.id}"
}

resource "aws_apigatewayv2_route" "delete" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "DELETE /items/{id}"
  target    = "integrations/${aws_apigatewayv2_integration.api.id}"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.api.id
  name        = "$default"
  auto_deploy = true

  tags = {
    ManagedBy = "aeroform"
  }
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.api.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.api.execution_arn}/*/*"
}
