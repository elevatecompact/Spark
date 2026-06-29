from pydantic_settings import BaseSettings, SettingsConfigDict


class Config(BaseSettings):
    model_config = SettingsConfigDict(env_file=None)

    MODERATION_PORT: int = 8102
    MODERATION_KAFKA_BROKERS: str = ""
    LOG_LEVEL: str = "INFO"
