variable "project_name" {
  description = "Name for the project — used in service name and labels"
  type        = string

  validation {
    condition     = can(regex("^[a-z0-9][a-z0-9-]{1,61}[a-z0-9]$", var.project_name))
    error_message = "Project name must be 3-63 characters, lowercase alphanumeric and hyphens only."
  }
}

variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region to deploy into"
  type        = string
  default     = "us-central1"
}

variable "container_image" {
  description = "Container image to deploy (e.g. gcr.io/my-project/my-app:latest)"
  type        = string
  default     = "gcr.io/cloudrun/hello"
}
