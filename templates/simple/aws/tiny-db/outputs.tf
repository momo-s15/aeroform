output "db_endpoint" {
  description = "Database connection endpoint (host:port)"
  value       = aws_db_instance.db.endpoint
}

output "db_host" {
  description = "Database hostname"
  value       = aws_db_instance.db.address
}

output "db_port" {
  description = "Database port"
  value       = aws_db_instance.db.port
}

output "db_name" {
  description = "Database name"
  value       = aws_db_instance.db.db_name
}

output "db_username" {
  description = "Database master username"
  value       = aws_db_instance.db.username
  sensitive   = true
}

output "connection_string" {
  description = "PostgreSQL connection string (add your password)"
  value       = "postgresql://${aws_db_instance.db.username}:PASSWORD@${aws_db_instance.db.endpoint}/${aws_db_instance.db.db_name}"
  sensitive   = true
}

output "security_note" {
  description = "This database is private — it can only be reached from within the VPC"
  value       = "Connect via an EC2 instance, Lambda, or VPN inside the same VPC."
}
