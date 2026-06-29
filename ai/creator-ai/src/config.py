import os


CREATOR_PORT: int = int(os.getenv("CREATOR_PORT", "8110"))
CREATOR_KAFKA_BROKERS: str = os.getenv("CREATOR_KAFKA_BROKERS", "")
LOG_LEVEL: str = os.getenv("LOG_LEVEL", "INFO")
