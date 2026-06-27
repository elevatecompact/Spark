# media-service — Event Contracts
## Published: media.upload.completed, media.transcoding.started, media.transcoding.completed, media.transcoding.failed, media.thumbnail.generated, media.drm.license.issued, media.content.deleted
## Consumed: stream.session.ended (process recording for VOD), commerce.order.fulfilled (grant media access), creator.channel.updated (process new avatar/banner), moderation.content.flagged (remove media), viewer.watch.started (CDN pre-warm for popular content)
## Schema: MediaTranscodingCompletedEvent {mediaId, sourceUrl, outputManifestUrl, profiles["720p","1080p","source"], duration, fileSize, completedAt}
