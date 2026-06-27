variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID for the EKS cluster"
  type        = string
}

variable "private_subnet_ids" {
  description = "List of private subnet IDs for the EKS cluster"
  type        = list(string)
}

variable "system_node_instance" {
  description = "Instance type for system node group"
  type        = string
  default     = "t3.medium"
}

variable "app_node_instance" {
  description = "Instance type for application node group"
  type        = string
  default     = "t3.large"
}

variable "min_nodes" {
  description = "Minimum number of nodes in the cluster"
  type        = number
  default     = 1
}

variable "max_nodes" {
  description = "Maximum number of nodes in the cluster"
  type        = number
  default     = 5
}

variable "desired_nodes" {
  description = "Desired number of nodes in the cluster"
  type        = number
  default     = 2
}

variable "cluster_version" {
  description = "Kubernetes version for the EKS cluster"
  type        = string
  default     = "1.30"
}

variable "region" {
  description = "AWS region"
  type        = string
}

variable "enable_karpenter" {
  description = "Enable Karpenter for node provisioning"
  type        = bool
  default     = false
}
