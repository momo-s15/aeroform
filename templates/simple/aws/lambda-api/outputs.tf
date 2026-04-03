output "api_url" {
  description = "Base URL for your API"
  value       = aws_apigatewayv2_stage.default.invoke_url
}

output "items_endpoint" {
  description = "Full endpoint for the items resource"
  value       = "${aws_apigatewayv2_stage.default.invoke_url}/items"
}

output "dynamodb_table_name" {
  description = "DynamoDB table name for your data"
  value       = aws_dynamodb_table.main.name
}

output "lambda_function_name" {
  description = "Lambda function name — check logs with: aws logs tail /aws/lambda/<name>"
  value       = aws_lambda_function.api.function_name
}

output "quick_test" {
  description = "Run this to test your API"
  value       = "curl -X POST ${aws_apigatewayv2_stage.default.invoke_url}/items -H 'Content-Type: application/json' -d '{\"name\":\"hello\"}'"
}
