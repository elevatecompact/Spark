# discovery-service — Testing Guide
## Unit: Feed construction logic (merge, deduplicate, diversity sample), trending score computation (velocity formula), category tree building (recursive hierarchy), editorial pick date validation.
## Integration: Feed rendering with all sources (recs + trending + editorial), category content pagination, collection CRUD with item ordering, A/B feed variant assignment.
## Load: 10000 feed requests/min with multi-source merge, trending recalculation for 5000 concurrent streams, category tree with 200 categories.
