import logging
import random

from src.models import (
    BoundingBox,
    Classification,
    ClassifyImageRequest,
    ClassifyImageResponse,
    DetectedObject,
    DetectFacesRequest,
    DetectFacesResponse,
    DetectObjectsRequest,
    DetectObjectsResponse,
    FaceInfo,
    OCRRequest,
    OCRResponse,
    OCRTextSegment,
)
from src.repository import VisionRepository

logger = logging.getLogger(__name__)

_MOCK_CLASSIFICATIONS: list[tuple[str, float, str | None]] = [
    ("person", 0.95, "people"),
    ("outdoor", 0.82, "scene"),
    ("nature", 0.76, "scene"),
    ("building", 0.65, "objects"),
    ("sky", 0.58, "scene"),
]


class VisionService:
    def __init__(self, repository: VisionRepository, event_producer) -> None:
        self._repository = repository
        self._event_producer = event_producer

    async def classify_image(self, request: ClassifyImageRequest) -> ClassifyImageResponse:
        predictions = [
            Classification(label=label, confidence=conf, category=cat)
            for label, conf, cat in _MOCK_CLASSIFICATIONS[: request.top_k]
        ]
        dominant_category = predictions[0].category if predictions else None
        response = ClassifyImageResponse(
            predictions=predictions, dominant_category=dominant_category
        )
        await self._repository.log_prediction(request.image_url, "classify", response.model_dump())
        await self._event_producer.publish(
            "vision.classify.completed",
            {"image_url": request.image_url, "dominant_category": dominant_category},
        )
        return response

    async def detect_objects(self, request: DetectObjectsRequest) -> DetectObjectsResponse:
        count = random.randint(2, 3)
        mock_labels = ["person", "car", "dog", "cat", "bottle", "chair"]
        objects: list[DetectedObject] = []
        for i in range(count):
            label = random.choice(mock_labels)
            box = BoundingBox(
                x=round(random.uniform(0, 0.8), 4),
                y=round(random.uniform(0, 0.8), 4),
                width=round(random.uniform(0.1, 0.5), 4),
                height=round(random.uniform(0.1, 0.5), 4),
            )
            objects.append(
                DetectedObject(
                    label=label,
                    confidence=round(random.uniform(request.confidence_threshold, 1.0), 4),
                    bounding_box=box,
                )
            )
        response = DetectObjectsResponse(
            objects=objects, object_count=count, model_version="noop-v1"
        )
        await self._repository.log_prediction(request.image_url, "detect", response.model_dump())
        await self._event_producer.publish(
            "vision.detect.completed",
            {"image_url": request.image_url, "object_count": count},
        )
        return response

    async def detect_faces(self, request: DetectFacesRequest) -> DetectFacesResponse:
        count = random.randint(0, 2)
        faces: list[FaceInfo] = []
        for _ in range(count):
            box = BoundingBox(
                x=round(random.uniform(0, 0.7), 4),
                y=round(random.uniform(0, 0.7), 4),
                width=round(random.uniform(0.1, 0.3), 4),
                height=round(random.uniform(0.1, 0.3), 4),
            )
            faces.append(
                FaceInfo(
                    bounding_box=box,
                    confidence=round(random.uniform(0.8, 1.0), 4),
                    age=random.choice(["20-30", "30-40", "40-50", None]),
                    gender=random.choice(["male", "female", None]),
                    emotions={
                        "happy": round(random.uniform(0, 1), 4),
                        "neutral": round(random.uniform(0, 1), 4),
                    }
                    if random.random() > 0.3
                    else None,
                )
            )
        response = DetectFacesResponse(
            faces=faces, face_count=count, model_version="noop-v1"
        )
        await self._repository.log_prediction(request.image_url, "faces", response.model_dump())
        await self._event_producer.publish(
            "vision.faces.completed",
            {"image_url": request.image_url, "face_count": count},
        )
        return response

    async def ocr(self, request: OCRRequest) -> OCRResponse:
        segments = [
            OCRTextSegment(
                text="Sample OCR text extracted from image.",
                bounding_box=BoundingBox(x=0.1, y=0.2, width=0.8, height=0.1),
                confidence=0.92,
            )
        ]
        response = OCRResponse(
            text="Sample OCR text extracted from image.",
            confidence=0.92,
            segments=segments,
            model_version="noop-v1",
        )
        await self._repository.log_prediction(request.image_url, "ocr", response.model_dump())
        await self._event_producer.publish(
            "vision.ocr.completed",
            {"image_url": request.image_url, "text_length": len(response.text)},
        )
        return response
