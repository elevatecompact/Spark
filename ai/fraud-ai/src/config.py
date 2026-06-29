from pydantic_settings import BaseSettings, SettingsConfigDict


class Config(BaseSettings):
    model_config = SettingsConfigDict(env_file=None)

    FRAUD_PORT: int = 8101
    FRAUD_KAFKA_BROKERS: str = ""
    LOG_LEVEL: str = "INFO"
