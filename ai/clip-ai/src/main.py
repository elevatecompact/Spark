import logging

import uvicorn
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from src.config import CLIP_PORT, LOG_LEVEL
from src.handler import router
from src.models import HealthResponse

logging.basicConfig(level=getattr(logging, LOG_LEVEL.upper(), logging.INFO))

app = FastAPI(title="Clip AI Service", version="0.1.0")

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
    uvicorn.run("src.main:app", host="0.0.0.0", port=CLIP_PORT, reload=True)
