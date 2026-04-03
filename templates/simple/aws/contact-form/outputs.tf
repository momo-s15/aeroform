output "cloudfront_url" {
  description = "The CloudFront distribution URL for your website"
  value       = "https://${aws_cloudfront_distribution.site.domain_name}"
}

output "s3_bucket_name" {
  description = "S3 bucket name — upload your site files here"
  value       = aws_s3_bucket.site.id
}

output "contact_api_url" {
  description = "POST to this URL to submit the contact form"
  value       = "${aws_apigatewayv2_stage.default.invoke_url}/contact"
}

output "upload_command" {
  description = "Run this to upload your site"
  value       = "aws s3 sync ./public s3://${aws_s3_bucket.site.id} --delete"
}

output "ses_verification_note" {
  description = "Check your email to verify the SES identity before the form can send"
  value       = "Verify ${var.contact_email} by clicking the link AWS sent to that address."
}
