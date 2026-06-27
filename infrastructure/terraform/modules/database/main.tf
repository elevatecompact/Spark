locals {
  name_prefix = "${var.project_name}-${var.environment}"
  db_name     = "spark_${var.environment}"
  master_username = "spark_admin"

  monitoring_interval = var.environment == "production" ? 10 : 60
  backup_window      = "03:00-04:00"
  maintenance_window = "sun:05:00-sun:06:00"
}

resource "random_password" "master" {
  length           = 24
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "aws_db_parameter_group" "postgres" {
  name        = "${local.name_prefix}-postgres-pg"
  family      = "postgres16"
  description = "Optimized parameter group for ${var.environment}"

  parameter {
    name         = "max_connections"
    value        = var.environment == "production" ? "500" : "100"
  }

  parameter {
    name         = "shared_buffers"
    value        = var.environment == "production" ? "{DBInstanceClassMemory/4}" : "{DBInstanceClassMemory/8}"
    apply_method = "pending-reboot"
  }

  parameter {
    name         = "effective_cache_size"
    value        = "{DBInstanceClassMemory*3/4}"
    apply_method = "pending-reboot"
  }

  parameter {
    name         = "work_mem"
    value        = var.environment == "production" ? "16384" : "4096"
  }

  parameter {
    name         = "maintenance_work_mem"
    value        = var.environment == "production" ? "2097152" : "65536"
  }

  parameter {
    name         = "random_page_cost"
    value        = "1.1"
  }

  parameter {
    name         = "log_min_duration_statement"
    value        = var.environment == "production" ? "1000" : "100"
  }

  parameter {
    name         = "log_connections"
    value        = "1"
  }

  parameter {
    name         = "log_disconnections"
    value        = "1"
  }

  tags = {
    Name        = "${local.name_prefix}-postgres-pg"
    Environment = var.environment
  }
}

resource "aws_db_subnet_group" "this" {
  name        = "${local.name_prefix}-db-subnet-group"
  subnet_ids  = var.database_subnet_ids
  description = "Database subnet group for ${var.environment}"

  tags = {
    Name        = "${local.name_prefix}-db-subnet-group"
    Environment = var.environment
  }
}

resource "aws_security_group" "rds" {
  name        = "${local.name_prefix}-rds-sg"
  description = "Security group for RDS PostgreSQL"
  vpc_id      = var.vpc_id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr]
    description = "PostgreSQL from VPC"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "All outbound traffic"
  }

  tags = {
    Name        = "${local.name_prefix}-rds-sg"
    Environment = var.environment
  }
}

resource "aws_db_instance" "primary" {
  identifier = "${local.name_prefix}-postgres"

  engine                       = "postgres"
  engine_version               = var.db_engine_version
  instance_class               = var.db_instance_class
  allocated_storage            = var.environment == "production" ? 200 : 50
  max_allocated_storage        = var.environment == "production" ? 1000 : 200
  storage_type                 = "gp3"
  storage_encrypted            = true
  kms_key_id                   = aws_kms_key.rds.arn

  db_name                      = local.db_name
  username                     = local.master_username
  password                     = random_password.master.result
  port                         = 5432

  db_subnet_group_name         = aws_db_subnet_group.this.name
  vpc_security_group_ids       = [aws_security_group.rds.id]
  parameter_group_name         = aws_db_parameter_group.postgres.name

  multi_az                     = var.multi_az
  publicly_accessible          = false
  deletion_protection          = var.deletion_protection
  skip_final_snapshot          = var.environment == "production" ? false : true
  final_snapshot_identifier    = var.environment == "production" ? "${local.name_prefix}-postgres-final" : null
  copy_tags_to_snapshot        = true

  backup_retention_period      = var.db_backup_retention
  backup_window                = local.backup_window
  maintenance_window           = local.maintenance_window

  monitoring_interval          = local.monitoring_interval
  monitoring_role_arn          = aws_iam_role.rds_monitoring.arn

  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]

  auto_minor_version_upgrade  = true
  performance_insights_enabled = true
  performance_insights_retention_period = var.environment == "production" ? 731 : 7

  tags = {
    Name        = "${local.name_prefix}-postgres"
    Environment = var.environment
  }
}

resource "aws_db_instance" "read_replica" {
  count = var.read_replica_count

  identifier = "${local.name_prefix}-postgres-replica-${count.index + 1}"

  engine                       = "postgres"
  engine_version               = var.db_engine_version
  instance_class               = var.read_replica_class != null ? var.read_replica_class : var.db_instance_class
  allocated_storage            = var.environment == "production" ? 200 : 50
  storage_type                 = "gp3"
  storage_encrypted            = true
  kms_key_id                   = aws_kms_key.rds.arn

  replicate_source_db          = aws_db_instance.primary.identifier

  vpc_security_group_ids       = [aws_security_group.rds.id]
  parameter_group_name         = aws_db_parameter_group.postgres.name

  publicly_accessible          = false
  copy_tags_to_snapshot        = true

  monitoring_interval          = local.monitoring_interval
  monitoring_role_arn          = aws_iam_role.rds_monitoring.arn

  performance_insights_enabled = true
  performance_insights_retention_period = var.environment == "production" ? 731 : 7

  tags = {
    Name        = "${local.name_prefix}-postgres-replica-${count.index + 1}"
    Environment = var.environment
  }
}

resource "aws_kms_key" "rds" {
  description             = "KMS key for RDS encryption - ${var.environment}"
  deletion_window_in_days = 30
  enable_key_rotation     = true

  tags = {
    Name        = "${local.name_prefix}-rds-kms"
    Environment = var.environment
  }
}

resource "aws_iam_role" "rds_monitoring" {
  name = "${local.name_prefix}-rds-monitoring-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "monitoring.rds.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name        = "${local.name_prefix}-rds-monitoring-role"
    Environment = var.environment
  }
}

resource "aws_iam_role_policy_attachment" "rds_monitoring" {
  role       = aws_iam_role.rds_monitoring.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonRDSEnhancedMonitoringRole"
}

resource "aws_secretsmanager_secret" "database" {
  name = "${local.name_prefix}-database-credentials"

  recovery_window_in_days = var.environment == "production" ? 30 : 0

  tags = {
    Name        = "${local.name_prefix}-database-credentials"
    Environment = var.environment
  }
}

resource "aws_secretsmanager_secret_version" "database" {
  secret_id = aws_secretsmanager_secret.database.id
  secret_string = jsonencode({
    username = local.master_username
    password = random_password.master.result
    host     = aws_db_instance.primary.address
    port     = aws_db_instance.primary.port
    database = local.db_name
    engine   = "postgres"
    engine_version = var.db_engine_version
  })
}

resource "aws_cloudwatch_metric_alarm" "cpu_high" {
  alarm_name          = "${local.name_prefix}-rds-cpu-high"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "CPUUtilization"
  namespace           = "AWS/RDS"
  period              = 300
  statistic           = "Average"
  threshold           = 80
  alarm_description   = "RDS CPU utilization is above 80%"
  alarm_actions       = []

  dimensions = {
    DBInstanceIdentifier = aws_db_instance.primary.identifier
  }

  tags = {
    Name        = "${local.name_prefix}-rds-cpu-high"
    Environment = var.environment
  }
}

resource "aws_cloudwatch_metric_alarm" "connections_high" {
  alarm_name          = "${local.name_prefix}-rds-connections-high"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "DatabaseConnections"
  namespace           = "AWS/RDS"
  period              = 300
  statistic           = "Average"
  threshold           = 80
  alarm_description   = "RDS database connections is above 80"
  alarm_actions       = []

  dimensions = {
    DBInstanceIdentifier = aws_db_instance.primary.identifier
  }

  tags = {
    Name        = "${local.name_prefix}-rds-connections-high"
    Environment = var.environment
  }
}

resource "aws_cloudwatch_metric_alarm" "storage_low" {
  alarm_name          = "${local.name_prefix}-rds-storage-low"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = 1
  metric_name         = "FreeStorageSpace"
  namespace           = "AWS/RDS"
  period              = 300
  statistic           = "Average"
  threshold           = 5000000000
  alarm_description   = "RDS free storage space is below 5GB"
  alarm_actions       = []

  dimensions = {
    DBInstanceIdentifier = aws_db_instance.primary.identifier
  }

  tags = {
    Name        = "${local.name_prefix}-rds-storage-low"
    Environment = var.environment
  }
}

resource "aws_cloudwatch_metric_alarm" "replica_lag" {
  count = var.read_replica_count > 0 ? 1 : 0

  alarm_name          = "${local.name_prefix}-rds-replica-lag"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "ReplicaLag"
  namespace           = "AWS/RDS"
  period              = 300
  statistic           = "Average"
  threshold           = 30
  alarm_description   = "RDS replica lag is above 30 seconds"
  alarm_actions       = []

  dimensions = {
    DBInstanceIdentifier = aws_db_instance.primary.identifier
  }

  tags = {
    Name        = "${local.name_prefix}-rds-replica-lag"
    Environment = var.environment
  }
}
