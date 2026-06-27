variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
}

variable "domain_name" {
  description = "Primary domain name"
  type        = string
}

variable "cloudfront_domain" {
  description = "CloudFront distribution domain name"
  type        = string
  default     = ""
}

variable "alb_domain" {
  description = "ALB DNS domain name"
  type        = string
  default     = ""
}

variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "alb_zone_id" {
  description = "Canonical hosted zone ID of the ALB"
  type        = string
  default     = ""
}

variable "dkim_selector" {
  description = "DKIM selector for SES"
  type        = string
  default     = "amazonses"
}

variable "create_zone" {
  description = "Create Route53 zone"
  type        = bool
  default     = true
}
