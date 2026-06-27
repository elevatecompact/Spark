output "media_bucket_id" {
  description = "ID of the media bucket"
  value       = aws_s3_bucket.media.id
}

output "media_bucket_name" {
  description = "Name of the media bucket"
  value       = aws_s3_bucket.media.bucket
}

output "media_bucket_arn" {
  description = "ARN of the media bucket"
  value       = aws_s3_bucket.media.arn
}

output "thumbnails_bucket_id" {
  description = "ID of the thumbnails bucket"
  value       = aws_s3_bucket.thumbnails.id
}

output "thumbnails_bucket_name" {
  description = "Name of the thumbnails bucket"
  value       = aws_s3_bucket.thumbnails.bucket
}

output "thumbnails_bucket_arn" {
  description = "ARN of the thumbnails bucket"
  value       = aws_s3_bucket.thumbnails.arn
}

output "logs_bucket_id" {
  description = "ID of the logs bucket"
  value       = aws_s3_bucket.logs.id
}

output "logs_bucket_name" {
  description = "Name of the logs bucket"
  value       = aws_s3_bucket.logs.bucket
}

output "logs_bucket_arn" {
  description = "ARN of the logs bucket"
  value       = aws_s3_bucket.logs.arn
}

output "backups_bucket_id" {
  description = "ID of the backups bucket"
  value       = aws_s3_bucket.backups.id
}

output "backups_bucket_name" {
  description = "Name of the backups bucket"
  value       = aws_s3_bucket.backups.bucket
}

output "backups_bucket_arn" {
  description = "ARN of the backups bucket"
  value       = aws_s3_bucket.backups.arn
}

output "terraform_state_bucket_id" {
  description = "ID of the terraform state bucket"
  value       = aws_s3_bucket.terraform_state.id
}

output "terraform_state_bucket_name" {
  description = "Name of the terraform state bucket"
  value       = aws_s3_bucket.terraform_state.bucket
}

output "terraform_state_bucket_arn" {
  description = "ARN of the terraform state bucket"
  value       = aws_s3_bucket.terraform_state.arn
}

output "cloudfront_logs_bucket_id" {
  description = "ID of the CloudFront logs bucket"
  value       = aws_s3_bucket.cloudfront_logs.id
}

output "cloudfront_logs_bucket_name" {
  description = "Name of the CloudFront logs bucket"
  value       = aws_s3_bucket.cloudfront_logs.bucket
}

output "cloudfront_logs_bucket_arn" {
  description = "ARN of the CloudFront logs bucket"
  value       = aws_s3_bucket.cloudfront_logs.arn
}

output "tempo_traces_bucket_id" {
  description = "ID of the Tempo traces bucket"
  value       = aws_s3_bucket.tempo_traces.id
}

output "tempo_traces_bucket_name" {
  description = "Name of the Tempo traces bucket"
  value       = aws_s3_bucket.tempo_traces.bucket
}

output "tempo_traces_bucket_arn" {
  description = "ARN of the Tempo traces bucket"
  value       = aws_s3_bucket.tempo_traces.arn
}

output "kms_key_arn" {
  description = "ARN of the KMS key used for S3 encryption"
  value       = aws_kms_key.storage.arn
}

output "kms_key_id" {
  description = "ID of the KMS key used for S3 encryption"
  value       = aws_kms_key.storage.key_id
}
