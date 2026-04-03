variable "name" {
  type        = string
  description = "Storage account name (3-24 lowercase alphanumeric, globally unique)"
}

variable "location" {
  type        = string
  description = "Azure region"
}

variable "resource_group_name" {
  type        = string
  description = "Resource group name"
}

variable "account_tier" {
  type        = string
  default     = "Standard"
  description = "Storage account tier (Standard or Premium)"
}

variable "replication_type" {
  type        = string
  default     = "LRS"
  description = "Replication type (LRS, GRS, RAGRS, ZRS)"
}
