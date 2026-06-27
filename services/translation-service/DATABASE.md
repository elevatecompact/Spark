# translation-service — Database Schema
## PostgreSQL — Translation memory & metadata
### 	ranslation_memory: source_hash VARCHAR(64) PK (SHA-256), source_text TEXT, translated_text TEXT, source_lang VARCHAR(10), target_lang VARCHAR(10), provider VARCHAR(20), quality_score FLOAT (0-1), created_at, updated_at. UNIQUE(source_hash, source_lang, target_lang)
### 	ranslation_jobs: id UUID PK, content_type VARCHAR, content_id UUID, status(pending,processing,completed,failed), languages TEXT[](targets), created_at
### eview_queue: id UUID PK, translation_id FK, original_text, translated_text, source_lang, target_lang, reviewer_id UUID nullable, status(pending,approved,rejected), reviewed_at
## Redis — Translation cache (TTL 24h), provider rate limit counters, language detection cache (TTL 1h)
