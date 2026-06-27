# notification-service — Database Schema
## 
otifications (inbox): id UUID PK, user_id FK, type VARCHAR, title, body TEXT, data JSONB (deep link payload), channel, read_at nullable, created_at. Retained 90 days.
## 
otification_preferences: user_id PK FK, preferences JSONB {type: {push:bool, email:bool, sms:bool, inapp:bool}}
## push_devices: id UUID PK, user_id FK, platform(ios,android,web), token (encrypted), is_active, created_at
## 	emplates: id UUID PK, type VARCHAR UNIQUE, subject_template, body_template (Handlebars), channels TEXT[] (which channels this template supports)
## Redis: Rate limit counters per channel per user, device token cache
