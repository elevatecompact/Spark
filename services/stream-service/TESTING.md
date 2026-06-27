# stream-service — Testing Guide
## Unit: Stream state machine transitions, ingest credential generation, transcoding profile resolution, HLS manifest generation.
## Integration: Full lifecycle (create→start→ingest→stop→archive), transcoding job submission, playback serving.
## Load: 50K concurrent viewers, 1000 simultaneous ingest connections, 100 transcoding jobs burst.
## Tools: docker-compose for local RTMP, OBS for test streams. k6 for load tests.
