# community-service — Troubleshooting
## Posts not appearing: Moderation filter caught it, community set to admin-only posting, post violated size limit. Check moderation flags, verify community posting permission, validate post length.
## Community not discoverable: Not featured, category not set, suspended. Check is_active, featured flag, verify category assignment.
## Member count wrong: Cache stale, count query slow, phantom members from failed leaves. Flush community cache, run count rebuild: ./community recount {communityId}, clean up stale memberships.
