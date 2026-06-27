variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
}

variable "grafana_admin_password" {
  description = "Admin password for Grafana"
  type        = string
  sensitive   = true
  default     = null
}

variable "loki_retention_days" {
  description = "Number of days to retain logs in Loki"
  type        = number
  default     = 7
}

variable "alertmanager_slack_webhook" {
  description = "Slack webhook URL for AlertManager notifications"
  type        = string
  sensitive   = true
  default     = null
}

variable "alertmanager_pagerduty_key" {
  description = "PagerDuty integration key for AlertManager"
  type        = string
  sensitive   = true
  default     = null
}

variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "tempo_retention_days" {
  description = "Number of days to retain traces in Tempo"
  type        = number
  default     = 3
}
