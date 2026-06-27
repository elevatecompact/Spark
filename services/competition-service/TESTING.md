# competition-service — Testing Guide
## Unit: Bracket generation algorithm (single/double elimination, bye handling), score validation, auto-progression logic, seeding fairness (random, ranked, manual), prize distribution calculation.
## Integration: Full competition lifecycle (create→register→start→matches→complete→prizes), match dispute workflow, judging panel scoring + aggregation, participant withdrawal mid-tournament.
## Load: Generate 256-participant bracket, 50 concurrent score submissions, 20 judge panel scoring.
