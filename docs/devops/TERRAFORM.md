# Terraform Infrastructure

## Approach

All Titan infrastructure is managed via Terraform using the official AWS provider. Environments are isolated through separate workspaces and state backends stored in S3 with DynamoDB locking.

## Module Structure

```
modules/
  networking/       VPC, subnets, NAT gateways, transit gateway
  eks/              EKS cluster, node groups, IRSA, add-ons
  rds/              PostgreSQL databases (Aurora, RDS)
  redis/            ElastiCache / MemoryDB for Redis
  s3/               S3 buckets with lifecycle policies
  iam/              IAM roles, policies, OIDC providers
  dns/              Route53 zones, records, health checks
  acm/              TLS certificates via AWS Certificate Manager
```

## State Management

Backend is `s3://titan-terraform-state-{env}` with DynamoDB table `titan-tf-locks`. All runs execute remotely via Terraform Cloud to ensure state consistency. Secrets are marked sensitive in output variables and never printed to logs.

## Workflow

Engineer opens a PR against `infra/live/{env}`. Atlantis auto-plans and posts the diff as a PR comment. After approval, `terraform apply` executes automatically. State is locked during apply, preventing concurrent mutations.

## Policy as Code

Sentinel or OPA policies enforce: mandatory tags on all resources, S3 buckets must have encryption and versioning enabled, security groups must not allow 0.0.0.0/0 on non-HTTP(S) ports, and EKS clusters must have private API endpoint enabled.