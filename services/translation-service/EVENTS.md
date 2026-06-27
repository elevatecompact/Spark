# translation-service — Event Contracts
## Published: translation.completed, translation.batch.completed, translation.review.ready, translation.provider.switched, translation.memory.updated
## Consumed: chat.message.sent (translate for multilingual rooms), stream.session.started (translate title/description), creator.channel.updated (translate new description), media.content.uploaded (translate metadata), community.post.created (translate posts)
## Schema: TranslationCompletedEvent {translationId, sourceText, translatedText, sourceLang, targetLang, provider(deepl|google), latencyMs, characterCount, timestamp}
