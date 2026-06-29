import logging
from contextlib import asynccontextmanager

import uvicorn
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from src.config import Config
from src.events import KafkaProducer
from src.handler import create_router
from src.models import HealthResponse
from src.repository import FraudRepository
from src.service import FraudService

logger = logging.getLogger(__name__)


def create_app() -> FastAPI:
    config = Config()
    logging.basicConfig(level=getattr(logging, config.LOG_LEVEL.upper(), logging.INFO))

    repo = FraudRepository()
    event_pub = KafkaProducer(brokers=config.FRAUD_KAFKA_BROKERS)
    service = FraudService(repo=repo, event_pub=event_pub)

    @asynccontextmanager
    async def lifespan(app: FastAPI):
        logger.info("Fraud AI Service starting on port %d", config.FRAUD_PORT)
        yield

    app = FastAPI(
        title="Fraud AI Service",
        version="0.1.0",
        lifespan=lifespan,
    )

    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    app.include_router(create_router(service))

    @app.get("/health", response_model=HealthResponse)
    async def health():
        return HealthResponse()

    return app


app = create_app()

if __name__ == "__main__":
    config = Config()
    uvicorn.run(
        "src.main:app",
        host="0.0.0.0",
        port=config.FRAUD_PORT,
        log_level=config.LOG_LEVEL.lower(),
    )
