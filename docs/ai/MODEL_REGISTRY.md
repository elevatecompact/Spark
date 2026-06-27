# Model Registry and Versioning

The Spark Model Registry is the central catalog of all ML models in production, staging, and development. Built on MLflow, it provides version control, lineage tracking, approval workflows, and artifact management.

## Registry Structure

Each model entry contains:
- **Metadata**: Model name, description, task type, framework, input/output schemas, license
- **Version History**: Every registered version with changelog, author, commit hash, and training dataset ID
- **Artifacts**: Model binary, preprocessing pipeline, configuration, and evaluation metrics
- **Lineage**: Links to training pipeline run, source code commit, dataset version, and feature set
- **Deployment Status**: Current stage (staging, production, archived, deprecated) with deployment history

## Version Lifecycle

| Stage | Description | Governance |
|---|---|---|
| Development | Experimental models, unvalidated | No restrictions |
| Staging | Validated on holdout set, waiting for shadow deployment | Peer review required |
| Production | Serving live traffic | Automated monitoring + human approval for promotion |
| Archived | Replaced by newer version | Retained for reproducibility |
| Deprecated | Known issues or retired | Scheduled for removal |

## Approval Workflow

1. Data scientist registers a new model version with evaluation metrics
2. Automated CI pipeline validates: reproducibility, fairness thresholds, latency benchmarks, and security scan
3. Peer review by another ML engineer
4. Model security review (input/output scanning for adversarial inputs)
5. Staging deployment and shadow evaluation (minimum 24 hours)
6. Production approval by ML lead or SRE
7. Canary deployment with automated monitoring

## Governance

All stage transitions are logged with timestamps, approver identity, and evidence artifacts. The registry supports audit queries for model lineage and training data provenance.
