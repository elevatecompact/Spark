# notification-service — Configuration
NOTIFICATION_PORT=4014, NOTIFICATION_DB_URL, NOTIFICATION_REDIS_URL, NOTIFICATION_KAFKA_BROKERS, PUSH_FCM_SERVER_KEY, PUSH_APNS_KEY_ID, PUSH_APNS_TEAM_ID, PUSH_APNS_KEY_PATH, SENDGRID_API_KEY, TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, TWILIO_PHONE_NUMBER, DIGEST_CRON="0 8 * * *" (daily), MAX_BATCH_SIZE=1000
FF: push_enabled=true, email_enabled=true, sms_enabled=false, inapp_enabled=true, digests_enabled=true
Rate limits: 50 push/min per device, 10 emails/min per user, 5 SMS/min per user
