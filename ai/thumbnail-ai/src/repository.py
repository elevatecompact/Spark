import logging

logger = logging.getLogger(__name__)


class ThumbnailRepository:
    def save_thumbnail(self, video_url: str, thumbnail_url: str, time_sec: float) -> None:
        logger.info("save_thumbnail(video_url=%s, thumbnail_url=%s, time_sec=%s)", video_url, thumbnail_url, time_sec)

    def log_extraction(self, video_url: str, frame_count: int) -> None:
        logger.info("log_extraction(video_url=%s, frame_count=%d)", video_url, frame_count)
