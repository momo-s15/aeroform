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

# --- Lambda function ---

resource "aws_iam_role" "lambda" {
  name = "${var.project_name}-cron-lambda-role"

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
    Name      = "${var.project_name}-cron-lambda-role"
    ManagedBy = "aeroform"
  }
}

resource "aws_iam_role_policy" "lambda" {
  name = "${var.project_name}-cron-lambda-policy"
  role = aws_iam_role.lambda.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["logs:CreateLogGroup", "logs:CreateLogStream", "logs:PutLogEvents"]
        Resource = "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:*"
      }
    ]
  })
}

resource "aws_lambda_function" "cron" {
  function_name = "${var.project_name}-cron"
  role          = aws_iam_role.lambda.arn
  handler       = "index.handler"
  runtime       = "nodejs20.x"
  timeout       = 60
  memory_size   = 128

  filename         = data.archive_file.lambda_placeholder.output_path
  source_code_hash = data.archive_file.lambda_placeholder.output_base64sha256

  environment {
    variables = {
      PROJECT_NAME = var.project_name
    }
  }

  tags = {
    Name      = "${var.project_name}-cron"
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
        console.log('Scheduled job running for', process.env.PROJECT_NAME);
        console.log('Event:', JSON.stringify(event));
        // Replace this with your scheduled task logic
        return { statusCode: 200, body: 'ok' };
      };
    JS
    filename = "index.js"
  }
}

resource "aws_cloudwatch_log_group" "cron" {
  name              = "/aws/lambda/${aws_lambda_function.cron.function_name}"
  retention_in_days = 7

  tags = { ManagedBy = "aeroform" }
}

# --- EventBridge scheduled rule ---

resource "aws_cloudwatch_event_rule" "cron" {
  name                = "${var.project_name}-cron-schedule"
  description         = "Scheduled trigger for ${var.project_name}"
  schedule_expression = var.schedule_expression

  tags = {
    Name      = "${var.project_name}-cron-schedule"
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

resource "aws_cloudwatch_event_target" "cron" {
  rule = aws_cloudwatch_event_rule.cron.name
  arn  = aws_lambda_function.cron.arn
}

resource "aws_lambda_permission" "eventbridge" {
  statement_id  = "AllowEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.cron.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.cron.arn
}
