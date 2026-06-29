import os


ASSISTANT_PORT: int = int(os.getenv("ASSISTANT_PORT", "8108"))
ASSISTANT_KAFKA_BROKERS: str = os.getenv("ASSISTANT_KAFKA_BROKERS", "")
LOG_LEVEL: str = os.getenv("LOG_LEVEL", "INFO")
