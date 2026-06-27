variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
}

variable "cross_region_replication" {
  description = "Enable cross-region replication for backup bucket"
  type        = bool
  default     = false
}

variable "replication_region" {
  description = "Destination region for cross-region replication"
  type        = string
  default     = "us-west-2"
}

variable "cloudfront_arn" {
  description = "ARN of the CloudFront distribution for bucket policy"
  type        = string
  default     = ""
}

variable "log_bucket_arn" {
  description = "ARN of the destination bucket for server access logs"
  type        = string
  default     = null
}
