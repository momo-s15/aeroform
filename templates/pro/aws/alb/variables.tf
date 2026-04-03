variable "name" {
  type        = string
  description = "ALB name"
}

variable "vpc_id" {
  type        = string
  description = "VPC ID for the target group"
}

variable "subnet_ids" {
  type        = list(string)
  description = "Public subnet IDs for the ALB"
}

variable "security_group_ids" {
  type        = list(string)
  description = "Security group IDs for the ALB"
}

variable "certificate_arn" {
  type        = string
  description = "ACM certificate ARN for HTTPS listener"
}

variable "target_port" {
  type        = number
  default     = 80
  description = "Port the target group forwards to"
}

variable "health_check_path" {
  type        = string
  default     = "/health"
  description = "Health check endpoint path"
}
