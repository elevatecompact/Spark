from src.models import ClipSegment


class ClipRepository:

    def save_clip_job(self, video_url: str, clips: list[ClipSegment]) -> None:
        pass

    def get_clip_jobs(self, limit: int = 10) -> list[dict]:
        return []
