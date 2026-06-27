# messaging-service — Testing Guide
## Unit: Conversation permissions, message validation, group size enforcement, reaction deduplication.
## Integration: Full lifecycle (create→message→read→delete), multi-device read state, attachment upload, group management, typing indicators.
## Load: 5000 conversations with 10K msgs each, 100 concurrent file uploads, 500-member group.
## Fixtures: ConversationFactory for randomized group setups.
