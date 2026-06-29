from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    THUMBNAIL_PORT: int = 8107
    THUMBNAIL_KAFKA_BROKERS: str = ""
    LOG_LEVEL: str = "INFO"

    model_config = {"env_prefix": ""}


settings = Settings()
