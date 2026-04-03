output "cloudfront_url" {
  description = "The CloudFront distribution URL for your website"
  value       = "https://${aws_cloudfront_distribution.site.domain_name}"
}

output "cloudfront_distribution_id" {
  description = "CloudFront distribution ID — use with aws cloudfront create-invalidation"
  value       = aws_cloudfront_distribution.site.id
}

output "s3_bucket_name" {
  description = "S3 bucket name — upload your site files here"
  value       = aws_s3_bucket.site.id
}

output "upload_command" {
  description = "Run this to upload your site"
  value       = "aws s3 sync ./public s3://${aws_s3_bucket.site.id} --delete"
}
