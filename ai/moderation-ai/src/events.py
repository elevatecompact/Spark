import logging
from abc import ABC, abstractmethod

logger = logging.getLogger(__name__)


class EventProducer(ABC):
    @abstractmethod
    async def publish(self, event_type: str, data: dict): ...


class NoopProducer(EventProducer):
    async def publish(self, event_type: str, data: dict):
        logger.info("Event [%s]: %s", event_type, data)


class KafkaProducer(EventProducer):
    def __init__(self, brokers: str):
        self._brokers = brokers

    async def publish(self, event_type: str, data: dict):
        if not self._brokers:
            logger.warning("No Kafka brokers configured, falling back to noop producer")
            await NoopProducer().publish(event_type, data)
            return
        logger.info("Kafka publish [%s] to %s: %s", event_type, self._brokers, data)
