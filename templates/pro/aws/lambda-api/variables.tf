variable "name" {
  type        = string
  description = "Lambda function and API Gateway name"
}

variable "memory_size" {
  type        = number
  default     = 256
  description = "Lambda memory in MB"
}

variable "timeout" {
  type        = number
  default     = 30
  description = "Lambda timeout in seconds"
}

variable "stage" {
  type        = string
  default     = "production"
  description = "Deployment stage name"
}
