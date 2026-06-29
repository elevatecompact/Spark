from pydantic_settings import BaseSettings, SettingsConfigDict


class Config(BaseSettings):
    model_config = SettingsConfigDict(env_file=None)

    TRANSLATION_PORT: int = 8103
    TRANSLATION_KAFKA_BROKERS: str = ""
    LOG_LEVEL: str = "INFO"
