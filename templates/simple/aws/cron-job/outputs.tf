output "lambda_function_name" {
  description = "Lambda function name — check logs with: aws logs tail /aws/lambda/<name>"
  value       = aws_lambda_function.cron.function_name
}

output "schedule" {
  description = "Current schedule expression"
  value       = var.schedule_expression
}

output "event_rule_name" {
  description = "EventBridge rule name"
  value       = aws_cloudwatch_event_rule.cron.name
}
