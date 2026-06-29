import uuid
from datetime import datetime, timezone

from src.events import EventProducer, create_event_producer
from src.models import (
    ChatMessage,
    ChatRequest,
    ChatResponse,
    ConversationHistory,
    GetHistoryResponse,
    ModifyContentRequest,
    ModifyContentResponse,
    SuggestRequest,
    SuggestResponse,
    SummarizeRequest,
    SummarizeResponse,
)
from src.repository import AssistantRepository


class AssistantService:
    def __init__(
        self,
        repository: AssistantRepository | None = None,
        event_producer: EventProducer | None = None,
    ) -> None:
        self.repository = repository or AssistantRepository()
        self.event_producer = event_producer or create_event_producer()

    async def chat(self, request: ChatRequest) -> ChatResponse:
        now = datetime.now(timezone.utc).isoformat()
        conversation_id = request.conversation_id or str(uuid.uuid4())

        conversation = self.repository.get_conversation(conversation_id)
        if conversation is None:
            conversation = ConversationHistory(
                conversation_id=conversation_id,
                messages=[],
                created_at=now,
                updated_at=now,
            )

        user_message = ChatMessage(role="user", content=request.message, timestamp=now)
        conversation.messages.append(user_message)

        reply_text = f"I received your message: '{request.message}'. This is a noop AI assistant response."
        assistant_message = ChatMessage(
            role="assistant", content=reply_text, timestamp=now
        )
        conversation.messages.append(assistant_message)

        tokens_used = len(request.message.split()) * 2 + 50
        conversation.updated_at = now

        self.repository.save_conversation(conversation)

        await self.event_producer.publish(
            "assistant.chat.completed",
            key=conversation_id,
            value={"conversation_id": conversation_id, "message": request.message},
        )

        return ChatResponse(
            reply=reply_text,
            conversation_id=conversation_id,
            tokens_used=tokens_used,
        )

    async def get_history(self, conversation_id: str) -> GetHistoryResponse:
        conversation = self.repository.get_conversation(conversation_id)
        if conversation is None:
            conversation = ConversationHistory(
                conversation_id=conversation_id,
                messages=[],
                created_at="",
                updated_at="",
            )
        return GetHistoryResponse(
            conversation=conversation,
            message_count=len(conversation.messages),
        )

    async def suggest(self, request: SuggestRequest) -> SuggestResponse:
        suggestions = [
            f"{request.context_type} suggestion {i + 1} for: {request.prompt}"
            for i in range(request.max_suggestions)
        ]
        await self.event_producer.publish(
            "assistant.suggest.completed",
            key=request.prompt,
            value={"context_type": request.context_type, "suggestions": suggestions},
        )
        return SuggestResponse(suggestions=suggestions, model_version="noop-v1")

    async def summarize(self, request: SummarizeRequest) -> SummarizeResponse:
        original_length = len(request.text)
        summary = request.text[: request.max_length]
        if len(request.text) > request.max_length:
            summary += "..."
        summary_length = len(summary)
        compression_ratio = (
            round(1 - (summary_length / original_length), 4)
            if original_length > 0
            else 0.0
        )

        await self.event_producer.publish(
            "assistant.summarize.completed",
            key=str(hash(request.text)),
            value={
                "original_length": original_length,
                "summary_length": summary_length,
            },
        )

        return SummarizeResponse(
            summary=summary,
            original_length=original_length,
            summary_length=summary_length,
            compression_ratio=compression_ratio,
            model_version="noop-v1",
        )

    async def modify_content(
        self, request: ModifyContentRequest
    ) -> ModifyContentResponse:
        instruction = request.instruction
        if request.tone:
            instruction += f" (tone: {request.tone})"
        if request.language:
            instruction += f" (language: {request.language})"

        modified_text = f"[Modified: {instruction}] {request.text}"
        changes_made = [f"Applied instruction: {request.instruction}"]
        if request.tone:
            changes_made.append(f"Adjusted tone to: {request.tone}")
        if request.language:
            changes_made.append(f"Translated to: {request.language}")

        await self.event_producer.publish(
            "assistant.modify.completed",
            key=str(hash(request.text)),
            value={"instruction": request.instruction},
        )

        return ModifyContentResponse(
            modified_text=modified_text,
            changes_made=changes_made,
            model_version="noop-v1",
        )
