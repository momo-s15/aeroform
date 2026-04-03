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

variable "db_name" {
  description = "Name of the database to create inside the instance"
  type        = string
  default     = "app"
}

variable "db_username" {
  description = "Master username for the database"
  type        = string
  default     = "aeroform"
}

variable "db_password" {
  description = "Master password for the database — use a strong value"
  type        = string
  sensitive   = true
}
