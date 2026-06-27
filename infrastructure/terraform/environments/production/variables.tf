variable "region" {
  description = "Primary AWS region"
  type        = string
  default     = "us-east-1"
}

variable "secondary_region" {
  description = "Secondary AWS region for disaster recovery"
  type        = string
  default     = "us-west-2"
}

variable "vpc_cidr" {
  description = "CIDR block for primary VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "secondary_vpc_cidr" {
  description = "CIDR block for secondary VPC"
  type        = string
  default     = "10.1.0.0/16"
}

variable "availability_zones" {
  description = "List of availability zones for primary region"
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b", "us-east-1c", "us-east-1d", "us-east-1e", "us-east-1f"]
}

variable "instance_types" {
  description = "EC2 instance types for app node group"
  type        = list(string)
  default     = ["c5.large", "c5.xlarge"]
}

variable "min_nodes" {
  description = "Minimum cluster nodes"
  type        = number
  default     = 5
}

variable "max_nodes" {
  description = "Maximum cluster nodes"
  type        = number
  default     = 20
}

variable "desired_nodes" {
  description = "Desired cluster nodes"
  type        = number
  default     = 8
}

variable "db_instance_class" {
  description = "RDS instance class for primary"
  type        = string
  default     = "db.r6g.xlarge"
}

variable "db_read_replica_class" {
  description = "RDS instance class for read replicas"
  type        = string
  default     = "db.r6g.large"
}

variable "db_read_replica_count" {
  description = "Number of read replicas"
  type        = number
  default     = 2
}

variable "db_backup_retention" {
  description = "Backup retention days"
  type        = number
  default     = 30
}

variable "db_multi_az" {
  description = "Enable multi-AZ deployment"
  type        = bool
  default     = true
}

variable "domain_name" {
  description = "Production domain name"
  type        = string
  default     = "sparkplatform.com"
}

variable "cloudflare_api_token" {
  description = "Cloudflare API token"
  type        = string
  sensitive   = true
}

variable "certificate_arn" {
  description = "ARN of existing ACM certificate (if pre-created)"
  type        = string
  default     = null
}
