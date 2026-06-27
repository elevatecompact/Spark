# licensing-service — Runbook
## Alerts: LicenseCheckLatency > 500ms, RoyaltyCalculationFailure, LicenseExpiryApproaching (30 days), PayoutDiscrepancyDetected, ComplianceViolationFlagged
## Process royalties manually: POST /v1/admin/royalties/process {periodStart, periodEnd}
## Terminate license: POST /v1/admin/licenses/{id}/terminate {reason}
## Override rights check: POST /v1/admin/rights/override {contentId, allowed:true}
## Generate compliance report: GET /v1/admin/compliance/report?period=2026-Q2
