# viewer-service — Testing Guide
## Unit: Progress calculation, preference validation, bookmark deduplication, rating range validation.
## Integration: Watch event batching, preference merge, bookmark CRUD, history pagination, cleanup batch job.
## Load: Simulate 10K concurrent viewers sending watch events. Measure bookmark list query with 5000 entries. k6: 	ests/performance/watch-events.js.
