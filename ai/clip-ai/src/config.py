from os import environ


CLIP_PORT: int = int(environ.get("CLIP_PORT", "8106"))
CLIP_KAFKA_BROKERS: str = environ.get("CLIP_KAFKA_BROKERS", "")
LOG_LEVEL: str = environ.get("LOG_LEVEL", "INFO")
