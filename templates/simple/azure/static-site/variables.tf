variable "project_name" {
  description = "Name for the project — used in resource names and tags"
  type        = string

  validation {
    condition     = can(regex("^[a-z0-9][a-z0-9-]{1,61}[a-z0-9]$", var.project_name))
    error_message = "Project name must be 3-63 characters, lowercase alphanumeric and hyphens only."
  }
}

variable "location" {
  description = "Azure region to deploy into (default suits many Azure for Education allow-lists; override via terraform.tfvars or AEROFORM_AZURE_LOCATION)"
  type        = string
  default     = "canadacentral"
}
