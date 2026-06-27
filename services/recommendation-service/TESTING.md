# recommendation-service — Testing Guide
## Unit: Feature computation logic, embedding similarity (cosine distance), feed deduplication, diversity sampling, business rule overlay (filter mature content).
## Integration: Full inference pipeline (feature→model→rank→serve), feedback loop incorporation, A/B test assignment.
## Evaluation: Offline metrics (NDCG@10, Recall@20, MAP) against held-out test set. Online CTR monitoring. Unit test model version compatibility.
## Tools: pytest for Python model code, go test for Go API server. Locust for load testing inference endpoint.
