from fastapi import APIRouter

from src.models import (
    ClassifyImageRequest,
    ClassifyImageResponse,
    DetectFacesRequest,
    DetectFacesResponse,
    DetectObjectsRequest,
    DetectObjectsResponse,
    OCRRequest,
    OCRResponse,
)
from src.service import VisionService

router = APIRouter(prefix="/v1/vision", tags=["vision"])


def _init(service: VisionService) -> None:
    router.service = service


@router.post("/classify", response_model=ClassifyImageResponse)
async def classify_image(request: ClassifyImageRequest) -> ClassifyImageResponse:
    return await router.service.classify_image(request)


@router.post("/detect", response_model=DetectObjectsResponse)
async def detect_objects(request: DetectObjectsRequest) -> DetectObjectsResponse:
    return await router.service.detect_objects(request)


@router.post("/faces", response_model=DetectFacesResponse)
async def detect_faces(request: DetectFacesRequest) -> DetectFacesResponse:
    return await router.service.detect_faces(request)


@router.post("/ocr", response_model=OCRResponse)
async def ocr(request: OCRRequest) -> OCRResponse:
    return await router.service.ocr(request)
