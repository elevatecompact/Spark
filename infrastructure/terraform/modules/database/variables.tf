variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID where the database will be deployed"
  type        = string
}

variable "database_subnet_ids" {
  description = "List of subnet IDs for the DB subnet group"
  type        = list(string)
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
  description = "Number of days to retain backups"
  type        = number
  default     = 7
}

variable "multi_az" {
  description = "Enable multi-AZ deployment"
  type        = bool
  default     = false
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
}

variable "deletion_protection" {
  description = "Enable deletion protection on RDS instance"
  type        = bool
  default     = true
}

variable "read_replica_count" {
  description = "Number of read replicas to create"
  type        = number
  default     = 0
}

variable "read_replica_class" {
  description = "Instance class for read replicas"
  type        = string
  default     = null
}
