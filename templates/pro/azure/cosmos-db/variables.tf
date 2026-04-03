variable "name" {
  type        = string
  description = "Cosmos DB account name (globally unique)"
}

variable "location" {
  type        = string
  description = "Azure region"
}

variable "resource_group_name" {
  type        = string
  description = "Resource group name"
}

variable "subnet_id" {
  type        = string
  description = "Subnet ID for the private endpoint"
}

variable "consistency_level" {
  type        = string
  default     = "Session"
  description = "Consistency level (BoundedStaleness, ConsistentPrefix, Eventual, Session, Strong)"
}
