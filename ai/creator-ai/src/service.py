import random
from datetime import datetime, timedelta, timezone

from src.events import EventProducer, NoopProducer
from src.models import (
    AnalyzeAudienceRequest,
    AnalyzeAudienceResponse,
    CaptionRequest,
    CaptionResponse,
    ContentIdea,
    ContentIdeaRequest,
    ContentIdeaResponse,
    HashtagRequest,
    HashtagResponse,
    HashtagSuggestion,
    OptimizeContentRequest,
    OptimizeContentResponse,
    ScheduleRequest,
    ScheduleResponse,
)
from src.repository import CreatorRepository


class CreatorService:

    def __init__(
        self,
        repository: CreatorRepository | None = None,
        event_producer: EventProducer | None = None,
    ):
        self._repository = repository or CreatorRepository()
        self._event_producer = event_producer or NoopProducer()

    async def generate_ideas(self, request: ContentIdeaRequest) -> ContentIdeaResponse:
        difficulties = ["easy", "medium", "hard"]
        engagements = ["low", "medium", "high"]
        formats = ["video", "post", "story", "reel", "carousel"]

        ideas = []
        for i in range(request.count):
            fmt = request.format or random.choice(formats)
            ideas.append(
                ContentIdea(
                    title=f"{request.niche.title()} Idea #{i + 1}",
                    description=f"A creative {fmt} idea about {request.niche} aimed at generating engagement.",
                    format=fmt,
                    difficulty=random.choice(difficulties),
                    estimated_engagement=random.choice(engagements),
                    tags=[request.niche.lower(), fmt, "content", "spark"],
                )
            )

        response = ContentIdeaResponse(
            ideas=ideas,
            count=len(ideas),
            model_version="noop-v1",
        )

        await self._event_producer.publish(
            "creator.ideas.completed",
            "ideas",
            {"count": request.count, "niche": request.niche},
        )

        return response

    async def optimize_content(self, request: OptimizeContentRequest) -> OptimizeContentResponse:
        response = OptimizeContentResponse(
            optimized_content=f"[Optimized] {request.content}",
            suggestions=[
                "Add a compelling hook in the first 3 seconds",
                "Use high-contrast visuals to improve retention",
                "Include a clear call-to-action at the end",
            ],
            score=round(random.uniform(5.0, 9.5), 1),
            improvements=[
                "Improved readability",
                "Stronger opening",
                "Better keyword placement",
            ],
            model_version="noop-v1",
        )

        await self._event_producer.publish(
            "creator.optimize.completed",
            "optimize",
            {"content_type": request.content_type, "platform": request.platform},
        )

        return response

    async def analyze_audience(self, request: AnalyzeAudienceRequest) -> AnalyzeAudienceResponse:
        followers = random.randint(10000, 500000)
        growth = round(random.uniform(2.5, 8.5), 2)
        engagement = round(random.uniform(1.5, 5.5), 2)

        response = AnalyzeAudienceResponse(
            total_followers=followers,
            growth_rate=growth,
            demographics={
                "age_groups": {"13-17": 15, "18-24": 35, "25-34": 30, "35-44": 12, "45+": 8},
                "genders": {"male": 45, "female": 50, "other": 5},
                "top_locations": ["United States", "India", "Brazil", "United Kingdom", "Canada"],
            },
            top_content_types=[
                {"type": "video", "percentage": 45},
                {"type": "image", "percentage": 30},
                {"type": "carousel", "percentage": 15},
                {"type": "text", "percentage": 10},
            ],
            best_posting_times=["08:00 UTC", "12:00 UTC", "18:00 UTC", "20:00 UTC"],
            engagement_rate=engagement,
            model_version="noop-v1",
        )

        await self._event_producer.publish(
            "creator.audience.completed",
            "audience",
            {"creator_id": request.creator_id, "timeframe": request.timeframe},
        )

        return response

    async def suggest_schedule(self, request: ScheduleRequest) -> ScheduleResponse:
        suggested = datetime.now(timezone.utc) + timedelta(days=3)
        suggested_time = suggested.replace(hour=14, minute=0, second=0, microsecond=0).isoformat()

        alt1 = suggested + timedelta(hours=2)
        alt2 = suggested - timedelta(hours=4)

        response = ScheduleResponse(
            suggested_time=suggested_time,
            alternatives=[
                alt1.replace(hour=16, minute=0, second=0, microsecond=0).isoformat(),
                alt2.replace(hour=10, minute=0, second=0, microsecond=0).isoformat(),
                alt1.replace(hour=20, minute=0, second=0, microsecond=0).isoformat(),
            ],
            reasoning="Based on your audience activity patterns, mid-afternoon slots show highest engagement for this content type.",
            model_version="noop-v1",
        )

        await self._event_producer.publish(
            "creator.schedule.completed",
            "schedule",
            {"content_type": request.content_type or "general"},
        )

        return response

    async def suggest_hashtags(self, request: HashtagRequest) -> HashtagResponse:
        base_tags = ["content", "creator", "trending", "viral", "spark", "new", "featured", "top", "daily", "mustsee"]
        selected = base_tags[: request.count]

        hashtags = [
            HashtagSuggestion(
                hashtag=f"#{tag}",
                popularity=round(random.uniform(0.1, 1.0), 2),
                relevance=round(random.uniform(0.3, 1.0), 2),
                category=request.category,
            )
            for tag in selected
        ]

        response = HashtagResponse(hashtags=hashtags, model_version="noop-v1")

        await self._event_producer.publish(
            "creator.hashtags.completed",
            "hashtags",
            {"count": request.count, "category": request.category},
        )

        return response

    async def generate_captions(self, request: CaptionRequest) -> CaptionResponse:
        captions = [
            f"Check out this amazing {request.image_description}! 🔥 Perfect for your daily inspiration. #{request.tone}",
            f"Can't get over this {request.image_description}. What do you think? Let us know in the comments!",
            f"Your daily dose of {request.image_description}. Save this for later and share with someone who needs to see it.",
        ]

        response = CaptionResponse(captions=captions, model_version="noop-v1")

        await self._event_producer.publish(
            "creator.captions.completed",
            "captions",
            {"tone": request.tone, "length": request.length},
        )

        return response
