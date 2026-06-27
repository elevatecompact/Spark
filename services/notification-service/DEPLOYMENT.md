# notification-service — Deployment Guide
K8s: k8s/notification-service/ — api (3x 512MB), push-worker (4x 1GB), email-worker (2x 1GB), sms-worker (1x 512MB). Workers scale on queue depth.
Deploy: ./notification migrate up, then kubectl apply -f k8s/notification-service/.
Channel workers consume from dedicated Kafka topics. Bounce handling: email bounces processed via SendGrid webhook, push via FCM/APNs feedback.
Health: /health (DB+Redis+Kafka+SendGrid+Twilio), /ready (workers registered), /metrics :4111.
