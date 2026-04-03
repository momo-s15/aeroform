variable "name" {
  type        = string
  description = "ECS cluster and service name"
}

variable "vpc_id" {
  type        = string
  description = "VPC ID"
}

variable "subnet_ids" {
  type        = list(string)
  description = "Private subnet IDs for ECS tasks"
}

variable "target_group_arn" {
  type        = string
  description = "ALB target group ARN"
}

variable "alb_security_group_ids" {
  type        = list(string)
  description = "Security group IDs of the ALB (for ingress rules)"
}

variable "container_image" {
  type        = string
  description = "Docker image URI (e.g. 123456789.dkr.ecr.us-east-1.amazonaws.com/app:latest)"
}

variable "container_port" {
  type        = number
  default     = 8080
  description = "Container port to expose"
}

variable "cpu" {
  type        = number
  default     = 256
  description = "Fargate task CPU units"
}

variable "memory" {
  type        = number
  default     = 512
  description = "Fargate task memory in MiB"
}

variable "desired_count" {
  type        = number
  default     = 2
  description = "Number of running tasks"
}
