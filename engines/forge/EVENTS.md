# Events

## Published Events
- job.created - Payload: { jobId, profile, inputUri, outputUri, createdAt }.
- job.progress - Payload: { jobId, progress, currentSegment, fps }.
- job.completed - Payload: { jobId, outputManifestUri, segmentsCount, totalDuration, averageFps }.
- job.failed - Payload: { jobId, stage, error, recoverable }.
- gpu.overloaded - Alert when GPU encoder utilization > 90%.
- gpu.recovered - GPU back to normal load.

## Subscribed Events
- transcode.profile.updated - Update encoding parameters for active jobs.
- transcode.priority.changed - Adjust job priority in queue.
