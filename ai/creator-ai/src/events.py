import abc
import json
import logging

logger = logging.getLogger(__name__)


class EventProducer(abc.ABC):

    @abc.abstractmethod
    async def publish(self, topic: str, key: str, value: dict) -> None:
        ...


class NoopProducer(EventProducer):

    async def publish(self, topic: str, key: str, value: dict) -> None:
        logger.debug("NoopProducer.publish(topic=%s, key=%s, value=%s)", topic, key, value)


class KafkaProducer(EventProducer):

    def __init__(self, brokers: str):
        self._brokers = brokers
        self._producer = None

    async def publish(self, topic: str, key: str, value: dict) -> None:
        if self._producer is None:
            from aiokafka import AIOKafkaProducer
            self._producer = AIOKafkaProducer(
                bootstrap_servers=self._brokers,
                value_serializer=lambda v: json.dumps(v).encode(),
            )
            await self._producer.start()
        await self._producer.send(topic, key=key.encode(), value=value)
