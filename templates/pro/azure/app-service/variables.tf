variable "name" {
  type        = string
  description = "App Service name (globally unique)"
}

variable "location" {
  type        = string
  description = "Azure region"
}

variable "resource_group_name" {
  type        = string
  description = "Resource group name"
}

variable "sku_name" {
  type        = string
  default     = "B1"
  description = "App Service Plan SKU (F1, B1, S1, P1v3, etc.)"
}

variable "node_version" {
  type        = string
  default     = "20-lts"
  description = "Node.js version for the application stack"
}

variable "app_settings" {
  type        = map(string)
  default     = {}
  description = "Application settings (environment variables)"
}
