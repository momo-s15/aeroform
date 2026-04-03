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

variable "db_password" {
  description = "Master password for the database — use a strong value"
  type        = string
  sensitive   = true
}

variable "container_image" {
  description = "Public container image for App Runner (e.g. public.ecr.aws/nginx/nginx:latest)"
  type        = string
  default     = "public.ecr.aws/nginx/nginx:latest"
}

variable "app_port" {
  description = "Port your application listens on inside the container"
  type        = string
  default     = "8080"
}
