# messaging-service — Configuration
MSG_PORT=4007, MSG_WS_PORT=4008, MSG_DB_URL, MSG_REDIS_URL, MAX_GROUP_SIZE=500, MESSAGE_MAX_LENGTH=4000, ATTACHMENT_MAX_SIZE_MB=100, E2EE_ENABLED=false, ATTACHMENT_STORAGE=s3://titan-msg-attachments
FF: e2ee_enabled=false, read_receipts=true, typing_indicators=true, message_editing=true, reactions_enabled=true
RL: 60 msgs/min per conversation, 50 conversations/day per user, 20 attachments/h, 1 reaction/msg per user
