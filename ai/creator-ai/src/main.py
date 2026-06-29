import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from src.config import CREATOR_PORT, LOG_LEVEL
from src.handler import router
from src.models import HealthResponse
from src.service import CreatorService

logging.basicConfig(level=getattr(logging, LOG_LEVEL.upper(), logging.INFO))


@asynccontextmanager
async def lifespan(app: FastAPI):
    app.state.service = CreatorService()
    yield


app = FastAPI(title="Creator AI Service", version="0.1.0", lifespan=lifespan)

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


app.include_router(router)


def main() -> None:
    import uvicorn
    uvicorn.run("src.main:app", host="0.0.0.0", port=CREATOR_PORT, reload=False)


if __name__ == "__main__":
    main()
