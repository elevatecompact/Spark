# moderation-service — Configuration
MODERATION_PORT=4018, MODERATION_DB_URL, MODERATION_REDIS_URL, MODERATION_KAFKA_BROKERS, ML_TEXT_CLASSIFIER_URL=http://ml-text:8501, ML_IMAGE_CLASSIFIER_URL=http://ml-image:8501, NSFW_THRESHOLD=0.85, TOXICITY_THRESHOLD=0.9, VIOLENCE_THRESHOLD=0.8, HATE_SPEECH_THRESHOLD=0.85, MAX_TEXT_SCAN_LENGTH=10000, MAX_IMAGE_SIZE_MB=10, AUTO_ACTION_ENABLED=true, HUMAN_REVIEW_THRESHOLD=0.6 (confidence below this → human review)
FF: auto_scan_text=true, auto_scan_images=true, realtime_chat_scan=true, stream_monitoring_enabled=false, human_review_enabled=true, appeals_enabled=true
Rate limits: 10 reports/h per user, 1000 scans/min (total), 10 appeals/day per user
