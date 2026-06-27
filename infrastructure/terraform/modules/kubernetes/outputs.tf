output "cluster_name" {
  description = "Name of the EKS cluster"
  value       = aws_eks_cluster.this.name
}

output "cluster_endpoint" {
  description = "Endpoint for the EKS cluster API"
  value       = aws_eks_cluster.this.endpoint
}

output "cluster_ca" {
  description = "Base64-encoded CA certificate for the cluster"
  value       = aws_eks_cluster.this.certificate_authority[0].data
}

output "cluster_arn" {
  description = "ARN of the EKS cluster"
  value       = aws_eks_cluster.this.arn
}

output "cluster_version" {
  description = "Kubernetes version of the cluster"
  value       = aws_eks_cluster.this.version
}

output "cluster_security_group_id" {
  description = "Security group ID attached to the EKS cluster"
  value       = aws_security_group.cluster.id
}

output "oidc_provider_arn" {
  description = "ARN of the OIDC provider"
  value       = aws_iam_openid_connect_provider.this.arn
}

output "oidc_provider_url" {
  description = "URL of the OIDC provider"
  value       = aws_iam_openid_connect_provider.this.url
}

output "node_role_system_arn" {
  description = "ARN of the system node group role"
  value       = aws_iam_role.system_node.arn
}

output "node_role_app_arn" {
  description = "ARN of the app node group role"
  value       = aws_iam_role.app_node.arn
}

output "cluster_autoscaler_role_arn" {
  description = "ARN of the cluster autoscaler IAM role"
  value       = aws_iam_role.cluster_autoscaler.arn
}

output "lb_controller_role_arn" {
  description = "ARN of the load balancer controller IAM role"
  value       = aws_iam_role.lb_controller.arn
}

output "external_dns_role_arn" {
  description = "ARN of the external DNS IAM role"
  value       = aws_iam_role.external_dns.arn
}

output "cert_manager_role_arn" {
  description = "ARN of the cert-manager IAM role"
  value       = aws_iam_role.cert_manager.arn
}

output "alb_domain" {
  description = "DNS domain for the ALB (to be filled after ingress creation)"
  value       = ""
}
