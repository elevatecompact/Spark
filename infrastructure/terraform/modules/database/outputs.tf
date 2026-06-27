output "arn" {
  description = "ARN of the RDS primary instance"
  value       = aws_db_instance.primary.arn
}

output "endpoint" {
  description = "RDS primary instance endpoint"
  value       = aws_db_instance.primary.endpoint
}

output "address" {
  description = "RDS primary instance address"
  value       = aws_db_instance.primary.address
}

output "port" {
  description = "RDS primary instance port"
  value       = aws_db_instance.primary.port
}

output "database_name" {
  description = "Database name"
  value       = aws_db_instance.primary.db_name
}

output "master_username" {
  description = "Master username for the database"
  value       = aws_db_instance.primary.username
}

output "secret_arn" {
  description = "ARN of the Secrets Manager secret containing database credentials"
  value       = aws_secretsmanager_secret.database.arn
}

output "kms_key_arn" {
  description = "ARN of the KMS key used for RDS encryption"
  value       = aws_kms_key.rds.arn
}

output "read_replica_endpoints" {
  description = "Endpoints of read replicas"
  value       = aws_db_instance.read_replica[*].endpoint
}

output "parameter_group_name" {
  description = "Name of the DB parameter group"
  value       = aws_db_parameter_group.postgres.name
}

output "subnet_group_name" {
  description = "Name of the DB subnet group"
  value       = aws_db_subnet_group.this.name
}

output "security_group_id" {
  description = "ID of the RDS security group"
  value       = aws_security_group.rds.id
}

output "cloudwatch_alarm_cpu_high" {
  description = "Name of the CPU high alarm"
  value       = aws_cloudwatch_metric_alarm.cpu_high.alarm_name
}

output "cloudwatch_alarm_storage_low" {
  description = "Name of the storage low alarm"
  value       = aws_cloudwatch_metric_alarm.storage_low.alarm_name
}
