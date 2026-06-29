import abc
import json
import logging

from src.config import RECOMMENDATION_KAFKA_BROKERS

logger = logging.getLogger(__name__)


class EventProducer(abc.ABC):
    @abc.abstractmethod
    async def start(self) -> None:
        ...

    @abc.abstractmethod
    async def stop(self) -> None:
        ...

    @abc.abstractmethod
    async def publish(self, topic: str, key: str, value: dict) -> None:
        ...


class NoopProducer(EventProducer):
    async def start(self) -> None:
        logger.debug("NoopProducer started")

    async def stop(self) -> None:
        logger.debug("NoopProducer stopped")

    async def publish(self, topic: str, key: str, value: dict) -> None:
        logger.debug("Event %s/%s: %s", topic, key, json.dumps(value))


class KafkaProducer(EventProducer):
    def __init__(self) -> None:
        self._producer = None

    async def start(self) -> None:
        if not RECOMMENDATION_KAFKA_BROKERS:
            logger.warning("Kafka brokers not configured, falling back to NoopProducer")
            return
        try:
            from aiokafka import AIOKafkaProducer

            self._producer = AIOKafkaProducer(
                bootstrap_servers=RECOMMENDATION_KAFKA_BROKERS,
            )
            await self._producer.start()
            logger.info("Kafka producer started on %s", RECOMMENDATION_KAFKA_BROKERS)
        except Exception:
            logger.exception("Failed to start Kafka producer")
            self._producer = None

    async def stop(self) -> None:
        if self._producer is not None:
            await self._producer.stop()
            logger.info("Kafka producer stopped")

    async def publish(self, topic: str, key: str, value: dict) -> None:
        if self._producer is None:
            logger.debug("Event %s/%s: %s (kafka unavailable)", topic, key, json.dumps(value))
            return
        try:
            await self._producer.send(
                topic=topic,
                key=key.encode(),
                value=json.dumps(value).encode(),
            )
        except Exception:
            logger.exception("Failed to publish event to %s/%s", topic, key)
