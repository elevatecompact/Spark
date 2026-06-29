import logging
from contextlib import asynccontextmanager

import uvicorn
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from src.config import LOG_LEVEL, RECOMMENDATION_PORT
from src.events import EventProducer, NoopProducer
from src.handler import router
from src.models import HealthResponse
from src.repository import RecommendationRepository
from src.service import RecommendationService

logging.basicConfig(level=getattr(logging, LOG_LEVEL.upper(), logging.INFO))
logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    repository = RecommendationRepository()
    producer: EventProducer = NoopProducer()
    await producer.start()
    service = RecommendationService(repository=repository, producer=producer)
    app.state.service = service
    logger.info("Recommendation AI Service started on port %d", RECOMMENDATION_PORT)
    yield
    if hasattr(app.state, "service"):
        await app.state.service._producer.stop()
    logger.info("Recommendation AI Service stopped")


app = FastAPI(title="Recommendation AI Service", version="0.1.0", lifespan=lifespan)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(router)


@app.get("/health", response_model=HealthResponse)
async def health() -> HealthResponse:
    return HealthResponse()


if __name__ == "__main__":
    uvicorn.run("src.main:app", host="0.0.0.0", port=RECOMMENDATION_PORT, reload=False)
