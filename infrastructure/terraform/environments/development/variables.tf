variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "List of availability zones"
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b"]
}

variable "instance_types" {
  description = "EC2 instance types for node groups"
  type        = list(string)
  default     = ["t3.medium"]
}

variable "min_nodes" {
  description = "Minimum cluster nodes"
  type        = number
  default     = 1
}

variable "max_nodes" {
  description = "Maximum cluster nodes"
  type        = number
  default     = 3
}

variable "desired_nodes" {
  description = "Desired cluster nodes"
  type        = number
  default     = 1
}

variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.medium"
}

variable "db_backup_retention" {
  description = "Backup retention days"
  type        = number
  default     = 3
}

variable "domain_name" {
  description = "Domain name for the platform"
  type        = string
  default     = "dev.sparkplatform.com"
}

variable "cloudflare_api_token" {
  description = "Cloudflare API token"
  type        = string
  sensitive   = true
}
