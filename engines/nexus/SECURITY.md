# Security

- Presigned URLs time-limited (1 hour) and per-asset scoped.
- File validation for mime type, virus scan, size limits.
- Bucket isolation per tenant with S3 bucket policies.
- SSE-S3 or SSE-KMS encryption at rest.
- Malware scanning with ClamAV; infected files quarantined.
- Deletion policies with grace periods and confirmation.
