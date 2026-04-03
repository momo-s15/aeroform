output "shorten_url" {
  description = "POST {url} to this endpoint to create a short link"
  value       = "${aws_apigatewayv2_stage.default.invoke_url}/shorten"
}

output "redirect_base_url" {
  description = "Short links redirect from this base URL"
  value       = aws_apigatewayv2_stage.default.invoke_url
}

output "dynamodb_table_name" {
  description = "DynamoDB table storing short URL mappings"
  value       = aws_dynamodb_table.urls.name
}
