# Security

- Input/output URI validation against allowlist patterns.
- Ephemeral scratch space per job with secure wipe on completion.
- Signed job tokens with permissions scoped to input/output buckets.
- Resource limits enforced via cgroups or job objects.
- Manifest integrity: HLS/DASH manifests signed with HMAC.
- Optional per-job encryption key for output segments.
