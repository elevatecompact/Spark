# community-service — Runbook
## Alerts: CommunityCreateErrorRate > 2%, PostLatencyP99 > 500ms, MemberJoinRateSpike (possible bot attack), ReportedPostBacklog > 100
## Suspend community: POST /v1/admin/communities/{id}/suspend
## Feature community: POST /v1/admin/communities/{id}/feature — promotes to discoverable.
## Remove member: DELETE /v1/admin/communities/{id}/members/{userId}
## Merge communities: Admin tool to merge duplicate communities.
