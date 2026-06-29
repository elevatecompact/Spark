from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    VISION_PORT: int = 8105
    VISION_KAFKA_BROKERS: str = ""
    LOG_LEVEL: str = "INFO"

    model_config = {"env_prefix": "", "case_sensitive": True}


settings = Settings()
