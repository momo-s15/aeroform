output "upload_api_url" {
  description = "POST to this URL with {filename, contentType} to get a presigned upload URL"
  value       = "${aws_apigatewayv2_stage.default.invoke_url}/upload"
}

output "download_api_url" {
  description = "GET this URL with ?key=filename to get a presigned download URL"
  value       = "${aws_apigatewayv2_stage.default.invoke_url}/download"
}

output "s3_bucket_name" {
  description = "S3 bucket where files are stored"
  value       = aws_s3_bucket.uploads.id
}
