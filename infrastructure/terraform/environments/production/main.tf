terraform {
  required_version = ">= 1.6"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
      configuration_alternatives = ["aws.secondary"]
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.20"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.9"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.5"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = ">= 4.0"
    }
  }

  backend "s3" {
    bucket         = "spark-terraform-state-production"
    key            = "production/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "spark-terraform-locks-production"
  }
}

provider "aws" {
  region = var.region

  default_tags {
    tags = {
      Project     = "spark"
      Environment = "production"
      ManagedBy   = "Terraform"
    }
  }
}

provider "aws" {
  alias  = "secondary"
  region = var.secondary_region

  default_tags {
    tags = {
      Project     = "spark"
      Environment = "production"
      ManagedBy   = "Terraform"
      Region      = var.secondary_region
    }
  }
}

module "networking" {
  source = "../../modules/networking"

  environment        = "production"
  project_name       = "spark"
  vpc_cidr           = var.vpc_cidr
  availability_zones = slice(var.availability_zones, 0, 3)
  single_nat_gateway = false
}

module "networking_secondary" {
  source = "../../modules/networking"

  environment        = "production-secondary"
  project_name       = "spark"
  vpc_cidr           = var.secondary_vpc_cidr
  availability_zones = ["${var.secondary_region}a", "${var.secondary_region}b", "${var.secondary_region}c"]
  single_nat_gateway = false
}

module "database" {
  source = "../../modules/database"

  environment         = "production"
  project_name        = "spark"
  vpc_id              = module.networking.vpc_id
  database_subnet_ids = module.networking.database_subnet_ids
  db_instance_class   = var.db_instance_class
  db_engine_version   = "16.3"
  db_backup_retention = var.db_backup_retention
  multi_az            = var.db_multi_az
  vpc_cidr            = module.networking.vpc_cidr
  deletion_protection = true
  read_replica_count  = var.db_read_replica_count
  read_replica_class  = var.db_read_replica_class
}

module "kubernetes" {
  source = "../../modules/kubernetes"

  environment          = "production"
  project_name         = "spark"
  vpc_id               = module.networking.vpc_id
  private_subnet_ids   = module.networking.private_subnet_ids
  region               = var.region
  system_node_instance = "t3.large"
  app_node_instance    = "c5.large"
  min_nodes            = var.min_nodes
  max_nodes            = var.max_nodes
  desired_nodes        = var.desired_nodes
  cluster_version      = "1.30"
  enable_karpenter     = true
}

module "monitoring" {
  source = "../../modules/monitoring"

  environment  = "production"
  project_name = "spark"
  region       = var.region
  loki_retention_days  = 30
  tempo_retention_days = 14
}

module "dns" {
  source = "../../modules/dns"

  environment       = "production"
  project_name      = "spark"
  domain_name       = var.domain_name
  region            = var.region
  cloudfront_domain = module.cdn.cloudfront_domain
  alb_domain        = module.kubernetes.alb_domain
  alb_zone_id       = ""
  create_zone       = true
}

module "cdn" {
  source = "../../modules/cdn"

  environment      = "production"
  project_name     = "spark"
  domain_name      = var.domain_name
  alb_domain       = module.kubernetes.alb_domain
  certificate_arn  = module.dns.certificate_arn
  s3_media_bucket  = module.storage.media_bucket_name
  price_class      = "PriceClass_All"
  geo_restriction_type = "none"
}

module "storage" {
  source = "../../modules/storage"

  environment              = "production"
  project_name             = "spark"
  cross_region_replication = true
  replication_region       = var.secondary_region
  cloudfront_arn           = module.cdn.cloudfront_arn
}

resource "aws_s3_bucket" "backup_replica" {
  provider = aws.secondary
  bucket   = "spark-production-backups-replica"

  tags = {
    Name        = "spark-production-backups-replica"
    Environment = "production"
  }
}

resource "aws_s3_bucket_versioning" "backup_replica" {
  provider = aws.secondary
  bucket   = aws_s3_bucket.backup_replica.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "backup_replica" {
  provider = aws.secondary
  bucket   = aws_s3_bucket.backup_replica.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
    bucket_key_enabled = true
  }
}

resource "aws_s3_bucket_public_access_block" "backup_replica" {
  provider = aws.secondary
  bucket   = aws_s3_bucket.backup_replica.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_lifecycle_configuration" "backup_replica" {
  provider = aws.secondary
  bucket   = aws_s3_bucket.backup_replica.id

  rule {
    id     = "transition-to-glacier-deep-archive"
    status = "Enabled"

    transition {
      days          = 30
      storage_class = "GLACIER_DEEP_ARCHIVE"
    }
  }
}

resource "aws_iam_role" "backup_replication" {
  name = "spark-production-s3-replication-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "s3.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name        = "spark-production-s3-replication-role"
    Environment = "production"
  }
}

data "aws_iam_policy_document" "backup_replication" {
  statement {
    actions = [
      "s3:GetReplicationConfiguration",
      "s3:ListBucket",
    ]
    resources = [module.storage.backups_bucket_arn]
  }

  statement {
    actions = [
      "s3:GetObjectVersionForReplication",
      "s3:GetObjectVersionAcl",
      "s3:GetObjectVersionTagging",
    ]
    resources = ["${module.storage.backups_bucket_arn}/*"]
  }

  statement {
    actions = [
      "s3:ReplicateObject",
      "s3:ReplicateDelete",
      "s3:ReplicateTags",
    ]
    resources = ["${aws_s3_bucket.backup_replica.arn}/*"]
  }
}

resource "aws_iam_policy" "backup_replication" {
  name   = "spark-production-s3-replication-policy"
  policy = data.aws_iam_policy_document.backup_replication.json

  tags = {
    Name        = "spark-production-s3-replication-policy"
    Environment = "production"
  }
}

resource "aws_iam_role_policy_attachment" "backup_replication" {
  role       = aws_iam_role.backup_replication.name
  policy_arn = aws_iam_policy.backup_replication.arn
}

resource "aws_s3_bucket_replication_configuration" "backups" {
  bucket = module.storage.backups_bucket_id
  role   = aws_iam_role.backup_replication.arn

  rule {
    id     = "cross-region-replication"
    status = "Enabled"

    destination {
      bucket        = aws_s3_bucket.backup_replica.arn
      storage_class = "STANDARD_IA"
    }
  }
}

resource "aws_kms_key" "backup_replica" {
  provider                = aws.secondary
  description             = "KMS key for backup replica encryption"
  deletion_window_in_days = 30
  enable_key_rotation     = true

  tags = {
    Name        = "spark-production-backup-replica-kms"
    Environment = "production"
  }
}

resource "aws_db_instance" "cross_region_replica" {
  provider = aws.secondary

  identifier = "spark-production-postgres-cross-region"

  engine                       = "postgres"
  engine_version               = "16.3"
  instance_class               = "db.r6g.large"
  allocated_storage            = 200
  max_allocated_storage        = 1000
  storage_type                 = "gp3"
  storage_encrypted            = true
  kms_key_id                   = aws_kms_key.backup_replica.arn

  replicate_source_db          = module.database.endpoint

  vpc_security_group_ids       = [aws_security_group.cross_region_db.id]
  publicly_accessible          = false
  deletion_protection          = true
  copy_tags_to_snapshot        = true

  backup_retention_period      = 7
  monitoring_interval          = 60
  monitoring_role_arn          = aws_iam_role.cross_region_monitoring[0].arn

  performance_insights_enabled = true
  performance_insights_retention_period = 7

  tags = {
    Name        = "spark-production-postgres-cross-region"
    Environment = "production"
  }
}

resource "aws_security_group" "cross_region_db" {
  provider    = aws.secondary
  name        = "spark-production-cross-region-db-sg"
  description = "Security group for cross-region DB replica"
  vpc_id      = module.networking_secondary.vpc_id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [var.secondary_vpc_cidr]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "spark-production-cross-region-db-sg"
    Environment = "production"
  }
}

resource "aws_iam_role" "cross_region_monitoring" {
  provider = aws.secondary
  name     = "spark-production-cross-region-monitoring"

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
    Name        = "spark-production-cross-region-monitoring"
    Environment = "production"
  }
}

resource "aws_iam_role_policy_attachment" "cross_region_monitoring" {
  provider   = aws.secondary
  role       = aws_iam_role.cross_region_monitoring.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonRDSEnhancedMonitoringRole"
}

resource "aws_iam_role" "backup" {
  name = "spark-production-backup-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "backup.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name        = "spark-production-backup-role"
    Environment = "production"
  }
}

resource "aws_iam_role_policy_attachment" "backup" {
  role       = aws_iam_role.backup.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSBackupServiceRolePolicyForBackup"
}

resource "aws_backup_vault" "primary" {
  name = "spark-production-backup-vault"

  tags = {
    Name        = "spark-production-backup-vault"
    Environment = "production"
  }
}

resource "aws_backup_plan" "daily" {
  name = "spark-production-daily-backup"

  rule {
    rule_name         = "daily-backup"
    target_vault_name = aws_backup_vault.primary.name
    schedule          = "cron(0 5 * * ? *)"
    start_window      = 60
    completion_window = 360
    recovery_point_tags = {
      Type        = "DailyBackup"
      Environment = "production"
    }

    lifecycle {
      delete_after = 30
    }

    copy_action {
      destination_vault_arn = aws_backup_vault.primary.arn
      lifecycle {
        delete_after = 90
      }
    }
  }

  rule {
    rule_name         = "weekly-backup"
    target_vault_name = aws_backup_vault.primary.name
    schedule          = "cron(0 6 ? * SUN *)"
    start_window      = 60
    completion_window = 360
    recovery_point_tags = {
      Type        = "WeeklyBackup"
      Environment = "production"
    }

    lifecycle {
      delete_after = 90
    }

    copy_action {
      destination_vault_arn = aws_backup_vault.primary.arn
      lifecycle {
        delete_after = 365
      }
    }
  }

  tags = {
    Name        = "spark-production-daily-backup"
    Environment = "production"
  }
}

resource "aws_backup_selection" "rds" {
  name         = "spark-production-rds-backup"
  plan_id      = aws_backup_plan.daily.id
  iam_role_arn = aws_iam_role.backup.arn

  resources = [
    module.database.arn
  ]
}

resource "aws_backup_selection" "s3" {
  name         = "spark-production-s3-backup"
  plan_id      = aws_backup_plan.daily.id
  iam_role_arn = aws_iam_role.backup.arn

  resources = [
    module.storage.backups_bucket_arn
  ]
}
