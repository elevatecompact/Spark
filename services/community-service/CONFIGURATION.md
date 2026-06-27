# community-service — Configuration
COMMUNITY_PORT=4019, COMMUNITY_DB_URL, COMMUNITY_REDIS_URL, COMMUNITY_KAFKA_BROKERS, MAX_COMMUNITIES_PER_CREATOR=5, MAX_COMMUNITIES_PER_USER=20, MAX_MEMBERS_PER_COMMUNITY=100000, POST_CHAR_LIMIT=10000, COMMENT_CHAR_LIMIT=2000, POSTS_PER_PAGE=25, COMMENTS_PER_PAGE=50
FF: communities_enabled=true, private_communities=true, post_reactions=true, threaded_comments=true, announcements_enabled=true
Rate limits: 5 posts/h per user, 20 comments/h, 10 communities/day per creator, 20 reactions/min per user
