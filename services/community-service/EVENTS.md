# community-service — Event Contracts
## Published: community.created, community.updated, community.deleted, community.member.joined, community.member.left, community.post.created, community.post.pinned, community.role.changed
## Consumed: creator.channel.created (auto-create fan community), identity.user.deleted (remove membership), moderation.content.flagged (hide post), subscription.activated (unlock private community), notification.push.sent (confirm community alert)
## Schema: CommunityMemberJoinedEvent {communityId, userId, joinedAt, memberCount, role}
