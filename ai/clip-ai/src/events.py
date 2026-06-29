import json
import logging
from abc import ABC, abstractmethod

from src.config import CLIP_KAFKA_BROKERS

logger = logging.getLogger(__name__)


class EventProducer(ABC):

    @abstractmethod
    async def start(self) -> None:
        ...

    @abstractmethod
    async def send(self, topic: str, key: str, value: dict) -> None:
        ...

    @abstractmethod
    async def stop(self) -> None:
        ...


class NoopProducer(EventProducer):

    async def start(self) -> None:
        logger.debug("NoopProducer started")

    async def send(self, topic: str, key: str, value: dict) -> None:
        logger.debug("NoopProducer send skipped topic=%s key=%s", topic, key)

    async def stop(self) -> None:
        logger.debug("NoopProducer stopped")


class KafkaProducer(EventProducer):

    def __init__(self) -> None:
        self._producer = None

    async def start(self) -> None:
        try:
            from aiokafka import AIOKafkaProducer
            self._producer = AIOKafkaProducer(
                bootstrap_servers=CLIP_KAFKA_BROKERS,
            )
            await self._producer.start()
        except Exception:
            logger.warning("Kafka unavailable, falling back to noop producer")
            self._producer = None

    async def send(self, topic: str, key: str, value: dict) -> None:
        if self._producer is None:
            return
        payload = json.dumps(value).encode()
        await self._producer.send(topic, key=key.encode(), value=payload)

    async def stop(self) -> None:
        if self._producer is not None:
            await self._producer.stop()


def create_producer() -> EventProducer:
    if CLIP_KAFKA_BROKERS:
        return KafkaProducer()
    return NoopProducer()
