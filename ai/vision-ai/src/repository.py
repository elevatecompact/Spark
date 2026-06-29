import logging

logger = logging.getLogger(__name__)


class VisionRepository:
    async def log_prediction(self, image_url: str, task: str, results: dict) -> None:
        logger.debug("Prediction logged: task=%s, image_url=%s", task, image_url)
