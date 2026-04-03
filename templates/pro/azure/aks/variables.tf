variable "name" {
  type        = string
  description = "AKS cluster name"
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
  description = "VNet subnet ID for the default node pool"
}

variable "kubernetes_version" {
  type        = string
  default     = "1.29"
  description = "Kubernetes version"
}

variable "vm_size" {
  type        = string
  default     = "Standard_D2s_v3"
  description = "VM size for the default node pool"
}

variable "node_count" {
  type        = number
  default     = 2
  description = "Initial node count"
}

variable "enable_auto_scaling" {
  type        = bool
  default     = true
  description = "Enable cluster autoscaler"
}

variable "min_count" {
  type        = number
  default     = 1
  description = "Minimum node count (when autoscaling enabled)"
}

variable "max_count" {
  type        = number
  default     = 5
  description = "Maximum node count (when autoscaling enabled)"
}

variable "private_cluster" {
  type        = bool
  default     = true
  description = "Enable private cluster (API server not publicly accessible)"
}
