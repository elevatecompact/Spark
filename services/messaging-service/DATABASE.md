# messaging-service — Database Schema
## conversations: id UUID PK, type(direct,group), name, icon_url, created_by FK, is_active
## conversation_members: conversation_id+user_id composite PK, role(admin,member), last_read_message_id, joined_at
## messages: id UUID PK, conversation_id FK, sender_id FK, content TEXT (encrypted if E2EE), content_type, reply_to UUID nullable, deleted_at (soft), created_at. Partitioned by conversation_id hash.
## Redis: Typing indicators (TTL 10s), online presence, unread counts per conversation
