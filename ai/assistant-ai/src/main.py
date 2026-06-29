import logging

import uvicorn
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from src.config import ASSISTANT_PORT, LOG_LEVEL
from src.handler import create_handler
from src.models import HealthResponse
from src.service import AssistantService

logging.basicConfig(level=getattr(logging, LOG_LEVEL.upper(), logging.INFO))
logger = logging.getLogger(__name__)

app = FastAPI(title="Assistant AI Service", version="0.1.0")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

service = AssistantService()
handler = create_handler(service)
app.include_router(handler)


@app.get("/health", response_model=HealthResponse)
async def health() -> HealthResponse:
    return HealthResponse()


def main() -> None:
    logger.info("Starting Assistant AI Service on port %d", ASSISTANT_PORT)
    uvicorn.run("src.main:app", host="0.0.0.0", port=ASSISTANT_PORT, reload=False)


if __name__ == "__main__":
    main()
