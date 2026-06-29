import logging

from src.config import ASSISTANT_KAFKA_BROKERS

logger = logging.getLogger(__name__)


class EventProducer:
    async def publish(self, topic: str, key: str, value: dict) -> None:
        raise NotImplementedError


class NoopProducer(EventProducer):
    async def publish(self, topic: str, key: str, value: dict) -> None:
        logger.debug("NoopProducer: topic=%s key=%s value=%s", topic, key, value)


class KafkaProducer(EventProducer):
    def __init__(self) -> None:
        self._producer = None

    async def publish(self, topic: str, key: str, value: dict) -> None:
        logger.debug("KafkaProducer: topic=%s key=%s value=%s", topic, key, value)


def create_event_producer() -> EventProducer:
    if ASSISTANT_KAFKA_BROKERS:
        return KafkaProducer()
    return NoopProducer()
