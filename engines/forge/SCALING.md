# Scaling

Forge scales by adding worker nodes each with one or more GPUs. A central Redis-backed job queue dispatches jobs to workers via pull model. Workers register on startup and request work when idle. Queue supports priority scheduling and fair-share across tenants. For live transcoding, a dedicated worker pool is reserved to avoid VOD job interference. Supports GPU partitioning (MIG on A100) for workload isolation.
