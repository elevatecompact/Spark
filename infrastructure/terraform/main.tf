provider "aws" {
  region = var.region

  default_tags {
    tags = {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "Terraform"
    }
  }
}

provider "kubernetes" {
  host                   = module.kubernetes.cluster_endpoint
  cluster_ca_certificate = base64decode(module.kubernetes.cluster_ca)
  exec {
    api_version = "client.authentication.k8s.io/v1beta1"
    command     = "aws"
    args        = ["eks", "get-token", "--cluster-name", module.kubernetes.cluster_name]
  }
}

provider "helm" {
  kubernetes {
    host                   = module.kubernetes.cluster_endpoint
    cluster_ca_certificate = base64decode(module.kubernetes.cluster_ca)
    exec {
      api_version = "client.authentication.k8s.io/v1beta1"
      command     = "aws"
      args        = ["eks", "get-token", "--cluster-name", module.kubernetes.cluster_name]
    }
  }
}

provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

module "networking" {
  source = "./modules/networking"

  environment       = var.environment
  project_name      = var.project_name
  vpc_cidr          = var.vpc_cidr
  availability_zones = var.availability_zones
  single_nat_gateway = var.environment != "production"
}

module "database" {
  source = "./modules/database"

  environment        = var.environment
  project_name       = var.project_name
  vpc_id             = module.networking.vpc_id
  database_subnet_ids = module.networking.database_subnet_ids
  db_instance_class  = var.db_instance_class
  db_engine_version  = var.db_engine_version
  db_backup_retention = var.db_backup_retention
  multi_az           = var.db_multi_az
  vpc_cidr           = module.networking.vpc_cidr
}

module "kubernetes" {
  source = "./modules/kubernetes"

  environment          = var.environment
  project_name         = var.project_name
  vpc_id               = module.networking.vpc_id
  private_subnet_ids   = module.networking.private_subnet_ids
  region               = var.region
  system_node_instance = var.system_node_instance
  app_node_instance    = var.app_node_instance
  min_nodes            = var.min_nodes
  max_nodes            = var.max_nodes
  desired_nodes        = var.desired_nodes
}

module "monitoring" {
  source = "./modules/monitoring"

  environment  = var.environment
  project_name = var.project_name
}

module "dns" {
  source = "./modules/dns"

  environment      = var.environment
  project_name     = var.project_name
  domain_name      = var.domain_name
  cloudfront_domain = module.cdn.cloudfront_domain
  alb_domain       = module.kubernetes.alb_domain
}

module "cdn" {
  source = "./modules/cdn"

  environment       = var.environment
  project_name      = var.project_name
  domain_name       = var.domain_name
  alb_domain        = module.kubernetes.alb_domain
  certificate_arn   = module.dns.certificate_arn
  s3_media_bucket   = module.storage.media_bucket_name
}

module "storage" {
  source = "./modules/storage"

  environment      = var.environment
  project_name     = var.project_name
  cross_region_replication = var.environment == "production"
}
