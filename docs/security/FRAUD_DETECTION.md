# Fraud Detection (fraud-ai)

Spark's fraud detection subsystem, **fraud-ai**, provides real-time and batch fraud analysis using machine learning models, rule-based systems, and behavioral analytics.

## Architecture

fraud-ai operates as a microservice that evaluates transactions, account activities, and API requests against fraud signals. It ingests events from the event bus, scores them in real time, and returns risk assessments to downstream services.

## Detection Methods

### Machine Learning Models

- **Anomaly detection** — Isolation Forest and Autoencoder-based models identify outliers in transaction patterns, login behavior, and data access
- **Classification models** — Gradient-boosted decision trees classify events as legitimate or fraudulent
- **Graph analysis** — Relationship graphs identify collusion networks, synthetic identities, and account farming

### Rule Engine

A configurable rules engine allows fraud analysts to define deterministic rules such as transaction value exceeds threshold for account tenure, account created less than 24 hours ago performing high-value actions, or multiple accounts accessing the platform from the same device fingerprint.

### Behavioral Baselines

Each user has a behavioral profile updated continuously: typical login times, geographic locations, devices, transaction amounts, and API endpoints accessed. Deviations increase the risk score.

## Response Actions

Based on the risk score, fraud-ai triggers:

- **Allow** — Risk below threshold, no action
- **Step-up** — Risk moderate, additional verification required (MFA, ID verification)
- **Block** — Risk high, transaction or action blocked
- **Review** — Flagged for manual review by the fraud team

## Feedback Loop

Fraud analysts review flagged events and provide feedback (true positive, false positive). This feedback is used to retrain models and tune rules, creating a continuous improvement cycle.
