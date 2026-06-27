# community-service — Database Schema
## communities: id UUID PK, name VARCHAR(100), description TEXT, creator_id FK, type(public,restricted,private), category VARCHAR(50), avatar_url, banner_url, rules TEXT[], member_count INT, post_count INT, is_active, created_at
## community_members: community_id+user_id PK, role(admin,moderator,member), joined_at, last_active_at
## community_posts: id UUID PK, community_id FK, author_id FK, title VARCHAR, content TEXT, is_pinned BOOLEAN, is_announcement BOOLEAN, reaction_counts JSONB {emoji: count}, comment_count INT, deleted_at (soft), created_at
## post_comments: id UUID PK, post_id FK, author_id FK, parent_id FK nullable (threaded), content TEXT, reaction_counts JSONB, deleted_at
## post_reactions: post_id+user_id+emoji PK
## Redis: Community member count cache, post view counters, active user tracking
