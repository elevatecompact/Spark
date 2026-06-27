# moderation-service — Testing Guide
## Unit: Rule evaluation engine (condition matching, severity assignment), action state machine, report deduplication, queue assignment algorithm (round-robin vs load-based), content type validation.
## Integration: Scan text→flag→auto-action pipeline, image scan→review queue→human decision→action, appeal workflow, rule CRUD with versioning.
## Accuracy tests: Compare auto-decisions vs human moderator decisions on labeled test set. Measure precision/recall per violation category.
## Load: k6 — scan-throughput.js (1000 text scans/s), image-scan-load.js (100 images/s).
