variable "name" {
  type        = string
  description = "GKE cluster name"
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
  description = "VPC network self-link or name"
}

variable "subnetwork" {
  type        = string
  description = "VPC subnetwork self-link or name"
}

variable "machine_type" {
  type        = string
  default     = "e2-medium"
  description = "Node pool machine type"
}

variable "node_count" {
  type        = number
  default     = 2
  description = "Initial node count per zone"
}

variable "min_node_count" {
  type        = number
  default     = 1
  description = "Minimum nodes for autoscaling"
}

variable "max_node_count" {
  type        = number
  default     = 5
  description = "Maximum nodes for autoscaling"
}

variable "private_endpoint" {
  type        = bool
  default     = false
  description = "Make the master endpoint private (no public access)"
}

variable "master_cidr" {
  type        = string
  default     = "172.16.0.0/28"
  description = "CIDR for the GKE master network"
}
