variable "name" {
  type        = string
  description = "GCS bucket name (globally unique)"
}

variable "project_id" {
  type        = string
  description = "GCP project ID"
}

variable "location" {
  type        = string
  default     = "US"
  description = "Bucket location (US, EU, or a specific region)"
}
