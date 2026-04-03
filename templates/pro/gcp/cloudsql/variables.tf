variable "name" {
  type        = string
  description = "Cloud SQL instance name"
}

variable "project_id" {
  type        = string
  description = "GCP project ID"
}

variable "region" {
  type        = string
  description = "GCP region"
}

variable "network" {
  type        = string
  description = "VPC network self-link for private IP"
}

variable "database_version" {
  type        = string
  default     = "POSTGRES_15"
  description = "Database version (POSTGRES_15, MYSQL_8_0, etc.)"
}

variable "tier" {
  type        = string
  default     = "db-f1-micro"
  description = "Machine tier (db-f1-micro, db-custom-2-8192, etc.)"
}

variable "disk_size" {
  type        = number
  default     = 10
  description = "Disk size in GB"
}

variable "high_availability" {
  type        = bool
  default     = false
  description = "Enable regional high availability"
}

variable "db_name" {
  type        = string
  default     = "app"
  description = "Default database name"
}

variable "db_user" {
  type        = string
  default     = "aeroform"
  description = "Default database user"
}

variable "db_password" {
  type        = string
  sensitive   = true
  description = "Database user password"
}
