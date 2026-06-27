# event-service — Database Schema
## events: id UUID PK, creator_id FK, title VARCHAR, description TEXT, category VARCHAR, type(virtual,inperson,hybrid), start_at TIMESTAMPTZ, end_at TIMESTAMPTZ, timezone VARCHAR, max_attendees INT, stream_id FK nullable, status(draft,published,cancelled,completed), cover_image_url, created_at
## event_ticket_tiers: id UUID PK, event_id FK, name VARCHAR, price_cents, quantity_total INT, quantity_sold INT, benefits TEXT[], sales_start_at, sales_end_at
## event_attendees: event_id+ticket_id+user_id PK, status(registered,attended,cancelled,no_show), registered_at, attended_at
## event_sessions: id UUID PK, event_id FK, title VARCHAR, speaker VARCHAR, start_at, end_at, stream_id FK nullable
## event_series: id UUID PK, creator_id FK, title, description, frequency(daily,weekly,monthly), day_of_week, start_time, timezone, next_event_at, is_active
## Redis: Event attendee count cache, ticket availability counter, event page view counter
