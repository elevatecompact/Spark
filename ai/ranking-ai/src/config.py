from pydantic_settings import BaseSettings, SettingsConfigDict


class Config(BaseSettings):
    model_config = SettingsConfigDict(env_file=None)

    RANKING_PORT: int = 8100
    RANKING_KAFKA_BROKERS: str = ""
    LOG_LEVEL: str = "INFO"
