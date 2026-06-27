variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.1.0.0/16"
}

variable "availability_zones" {
  description = "List of availability zones"
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b", "us-east-1c"]
}

variable "instance_types" {
  description = "EC2 instance types for node groups"
  type        = list(string)
  default     = ["t3.large"]
}

variable "min_nodes" {
  description = "Minimum cluster nodes"
  type        = number
  default     = 2
}

variable "max_nodes" {
  description = "Maximum cluster nodes"
  type        = number
  default     = 5
}

variable "desired_nodes" {
  description = "Desired cluster nodes"
  type        = number
  default     = 3
}

variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.large"
}

variable "db_backup_retention" {
  description = "Backup retention days"
  type        = number
  default     = 14
}

variable "domain_name" {
  description = "Domain name for the platform"
  type        = string
  default     = "staging.sparkplatform.com"
}

variable "cloudflare_api_token" {
  description = "Cloudflare API token"
  type        = string
  sensitive   = true
}
