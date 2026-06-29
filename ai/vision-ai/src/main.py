import logging

import uvicorn
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from src.config import settings
from src.events import NoopProducer
from src.handler import _init as init_handler
from src.handler import router
from src.models import HealthResponse
from src.repository import VisionRepository
from src.service import VisionService

logging.basicConfig(level=getattr(logging, settings.LOG_LEVEL.upper(), logging.INFO))
logger = logging.getLogger(__name__)

app = FastAPI(title="Vision AI Service")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/health", response_model=HealthResponse)
async def health() -> HealthResponse:
    return HealthResponse()


@app.on_event("startup")
async def startup() -> None:
    repository = VisionRepository()
    producer = NoopProducer()
    service = VisionService(repository=repository, event_producer=producer)
    init_handler(service)
    app.include_router(router)
    logger.info("Vision AI service started")


if __name__ == "__main__":
    uvicorn.run("src.main:app", host="0.0.0.0", port=settings.VISION_PORT, reload=False)
