from pydantic import BaseModel


class BoundingBox(BaseModel):
    x: float
    y: float
    width: float
    height: float


class Classification(BaseModel):
    label: str
    confidence: float
    category: str | None


class ClassifyImageRequest(BaseModel):
    image_url: str
    categories: list[str] | None = None
    top_k: int = 5


class ClassifyImageResponse(BaseModel):
    predictions: list[Classification]
    dominant_category: str | None
    model_version: str = "noop-v1"


class DetectedObject(BaseModel):
    label: str
    confidence: float
    bounding_box: BoundingBox | None


class DetectObjectsRequest(BaseModel):
    image_url: str
    confidence_threshold: float = 0.5


class DetectObjectsResponse(BaseModel):
    objects: list[DetectedObject]
    object_count: int
    model_version: str


class FaceInfo(BaseModel):
    bounding_box: BoundingBox
    confidence: float
    age: str | None
    gender: str | None
    emotions: dict[str, float] | None


class DetectFacesRequest(BaseModel):
    image_url: str


class DetectFacesResponse(BaseModel):
    faces: list[FaceInfo]
    face_count: int
    model_version: str


class OCRTextSegment(BaseModel):
    text: str
    bounding_box: BoundingBox
    confidence: float


class OCRRequest(BaseModel):
    image_url: str
    language: str | None = None


class OCRResponse(BaseModel):
    text: str
    confidence: float
    segments: list[OCRTextSegment] | None
    model_version: str


class HealthResponse(BaseModel):
    status: str = "ok"
    service: str = "vision-ai"
    version: str = "0.1.0"
