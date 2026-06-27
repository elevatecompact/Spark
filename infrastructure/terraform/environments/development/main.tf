terraform {
  required_version = ">= 1.6"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
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
    bucket         = "spark-terraform-state-development"
    key            = "development/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "spark-terraform-locks-development"
  }
}

provider "aws" {
  region = var.region

  default_tags {
    tags = {
      Project     = "spark"
      Environment = "development"
      ManagedBy   = "Terraform"
    }
  }
}

module "networking" {
  source = "../../modules/networking"

  environment        = "development"
  project_name       = "spark"
  vpc_cidr           = var.vpc_cidr
  availability_zones = var.availability_zones
  single_nat_gateway = true
}

module "database" {
  source = "../../modules/database"

  environment         = "development"
  project_name        = "spark"
  vpc_id              = module.networking.vpc_id
  database_subnet_ids = module.networking.database_subnet_ids
  db_instance_class   = var.db_instance_class
  db_engine_version   = "16.3"
  db_backup_retention = 3
  multi_az            = false
  vpc_cidr            = module.networking.vpc_cidr
  deletion_protection = false
  read_replica_count  = 0
}

module "kubernetes" {
  source = "../../modules/kubernetes"

  environment          = "development"
  project_name         = "spark"
  vpc_id               = module.networking.vpc_id
  private_subnet_ids   = module.networking.private_subnet_ids
  region               = var.region
  system_node_instance = "t3.medium"
  app_node_instance    = "t3.medium"
  min_nodes            = var.min_nodes
  max_nodes            = var.max_nodes
  desired_nodes        = var.desired_nodes
  cluster_version      = "1.30"
}

module "monitoring" {
  source = "../../modules/monitoring"

  environment  = "development"
  project_name = "spark"
  region       = var.region
}

module "dns" {
  source = "../../modules/dns"

  environment       = "development"
  project_name      = "spark"
  domain_name       = var.domain_name
  region            = var.region
  cloudfront_domain = module.cdn.cloudfront_domain
  alb_domain        = ""
  alb_zone_id       = ""
}

module "cdn" {
  source = "../../modules/cdn"

  environment      = "development"
  project_name     = "spark"
  domain_name      = var.domain_name
  alb_domain       = module.kubernetes.alb_domain
  certificate_arn  = module.dns.certificate_arn
  s3_media_bucket  = module.storage.media_bucket_name
  price_class      = "PriceClass_100"
  geo_restriction_type = "none"
}

module "storage" {
  source = "../../modules/storage"

  environment              = "development"
  project_name             = "spark"
  cross_region_replication = false
  cloudfront_arn           = module.cdn.cloudfront_arn
}
