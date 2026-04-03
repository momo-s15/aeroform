output "app_url" {
  description = "Your application URL (HTTPS provided by App Runner)"
  value       = "https://${aws_apprunner_service.app.service_url}"
}

output "s3_bucket_name" {
  description = "S3 bucket for static assets"
  value       = aws_s3_bucket.assets.id
}

output "db_endpoint" {
  description = "Database connection endpoint (host:port)"
  value       = aws_db_instance.db.endpoint
}

output "db_connection_string" {
  description = "PostgreSQL connection string (password omitted)"
  value       = "postgresql://aeroform:PASSWORD@${aws_db_instance.db.endpoint}/app"
  sensitive   = true
}

output "app_runner_service_id" {
  description = "App Runner service ID"
  value       = aws_apprunner_service.app.service_id
}
