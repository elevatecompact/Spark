import logging
from abc import ABC, abstractmethod

logger = logging.getLogger(__name__)


class EventProducer(ABC):
    @abstractmethod
    async def publish(self, topic: str, key: str, payload: dict) -> None: ...


class NoopProducer(EventProducer):
    async def publish(self, topic: str, key: str, payload: dict) -> None:
        logger.debug("noop publish topic=%s key=%s payload=%s", topic, key, payload)


class KafkaProducer(EventProducer):
    def __init__(self, brokers: str):
        self._brokers = brokers
        self._producer = None

    async def publish(self, topic: str, key: str, payload: dict) -> None:
        logger.debug("kafka publish topic=%s key=%s payload=%s", topic, key, payload)
