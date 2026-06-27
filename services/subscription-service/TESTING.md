# subscription-service — Testing Guide
## Unit: Billing period calc, plan change proration, grace period eligibility, cancellation logic, invoice generation.
## Integration: Full lifecycle (subscribe→bill→cancel→expire), failed payment retry, plan upgrade proration, multi-sub per user.
## Billing cycle: Simulate 10K subs due for renewal, payment retry with partial success/failure, concurrent billing edge cases.
