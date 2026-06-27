# Data Retention Policies

SPARK maintains data retention policies that balance operational requirements with privacy regulations including GDPR, CCPA, and other applicable laws.

## Retention Schedules

**User Data.** Active user accounts are retained indefinitely. Deactivated accounts are soft-deleted and purged after 90 days. Users may request immediate deletion through the privacy center. Anonymized analytics data derived from user activity is retained separately without personal identifiers.

**Content Data.** Published content is retained for the life of the platform. Deleted content is soft-deleted and permanently purged after 30 days. Content metadata and engagement metrics are retained in aggregated form after content deletion.

**Transaction Data.** Financial transaction records are retained for 7 years to comply with tax and accounting regulations. After 7 years, personally identifiable information is anonymized while aggregate financial data is preserved.

**Event and Log Data.** Raw event data in ClickHouse is retained for 90 days. Pre-aggregated rollups are retained for 24 months. Application logs are retained for 30 days. Audit logs are retained for 3 years.

**Session and Cache Data.** Session data in Redis expires after 24 hours. Cache entries have configurable TTLs from 60 seconds to 24 hours depending on content type.

## Enforcement

Retention policies are enforced through automated purge jobs running as Kubernetes CronJobs. Each job logs the number of records purged and reports metrics to the monitoring system. Failed purge jobs trigger alerts for manual investigation.
