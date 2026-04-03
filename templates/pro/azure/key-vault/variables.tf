variable "name" {
  type        = string
  description = "Key Vault name (globally unique, 3-24 alphanumeric + hyphens)"
}

variable "location" {
  type        = string
  description = "Azure region"
}

variable "resource_group_name" {
  type        = string
  description = "Resource group name"
}
