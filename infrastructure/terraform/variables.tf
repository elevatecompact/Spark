variable "environment" {
  description = "Deployment environment (development, staging, production)"
  type        = string
  validation {
    condition     = contains(["development", "staging", "production"], var.environment)
    error_message = "Environment must be one of: development, staging, production."
  }
}

variable "region" {
  description = "AWS region for resources"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Name of the Spark project"
  type        = string
  default     = "spark"
}

variable "domain_name" {
  description = "Primary domain name for the platform"
  type        = string
  default     = "sparkplatform.com"
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "List of availability zones to use"
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b", "us-east-1c"]
}

variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.medium"
}

variable "db_engine_version" {
  description = "PostgreSQL engine version"
  type        = string
  default     = "16.3"
}

variable "db_backup_retention" {
  description = "Number of days to retain database backups"
  type        = number
  default     = 7
}

variable "db_multi_az" {
  description = "Enable multi-AZ deployment for RDS"
  type        = bool
  default     = false
}

variable "system_node_instance" {
  description = "EC2 instance type for system node group"
  type        = string
  default     = "t3.medium"
}

variable "app_node_instance" {
  description = "EC2 instance type for application node group"
  type        = string
  default     = "t3.large"
}

variable "min_nodes" {
  description = "Minimum number of cluster nodes"
  type        = number
  default     = 1
}

variable "max_nodes" {
  description = "Maximum number of cluster nodes"
  type        = number
  default     = 5
}

variable "desired_nodes" {
  description = "Desired number of cluster nodes"
  type        = number
  default     = 2
}

variable "cloudflare_api_token" {
  description = "Cloudflare API token for DNS management"
  type        = string
  sensitive   = true
}
