import abc
import json
import logging

logger = logging.getLogger(__name__)


class EventProducer(abc.ABC):
    @abc.abstractmethod
    async def publish(self, topic: str, payload: dict) -> None: ...


class NoopProducer(EventProducer):
    async def publish(self, topic: str, payload: dict) -> None:
        logger.debug("Event (noop): topic=%s, payload=%s", topic, json.dumps(payload))


class KafkaProducer(EventProducer):
    async def publish(self, topic: str, payload: dict) -> None:
        logger.info("Event (kafka): topic=%s, payload=%s", topic, json.dumps(payload))
