variable "name" {
  type        = string
  description = "CloudFront distribution name"
}

variable "origin_domain_name" {
  type        = string
  description = "Origin domain (ALB DNS name or API Gateway endpoint)"
}

variable "acm_certificate_arn" {
  type        = string
  default     = ""
  description = "ACM certificate ARN for custom domain (leave empty for CloudFront default)"
}

variable "enable_waf" {
  type        = bool
  default     = true
  description = "Attach a WAFv2 web ACL with rate limiting"
}

variable "waf_rate_limit" {
  type        = number
  default     = 2000
  description = "Maximum requests per 5-minute window per IP"
}
