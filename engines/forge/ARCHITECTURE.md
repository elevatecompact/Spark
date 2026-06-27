# Architecture

Forge uses a pipeline-of-stages architecture. Input demuxing runs in a dedicated thread, feeding raw frames into a frame pool. The distributor fans frames out to parallel encoding pipelines - one per rung in the bitrate ladder. Each encoder stage is backed by hardware encoding (NVENC or VAAPI) with software fallback. Encoded packets flow into a segmenter producing fragmented MP4s, then into a packager generating HLS and DASH manifests. A watchdog monitors health and triggers keyframe-request recovery on encoder stalls.
