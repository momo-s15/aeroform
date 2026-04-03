variable "project_name" {
  description = "Name for the project — used in resource names and tags"
  type        = string

  validation {
    condition     = can(regex("^[a-z0-9][a-z0-9-]{1,61}[a-z0-9]$", var.project_name))
    error_message = "Project name must be 3-63 characters, lowercase alphanumeric and hyphens only."
  }
}

variable "region" {
  description = "AWS region to deploy into"
  type        = string
  default     = "us-east-1"
}

variable "contact_email" {
  description = "Email address that receives contact form submissions (must be verified in SES)"
  type        = string
}
