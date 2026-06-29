package domain

import (
	"time"

	"github.com/google/uuid"
)

type ContentType string

const (
	MsgText     ContentType = "text"
	MsgImage    ContentType = "image"
	MsgVoice    ContentType = "voice"
	MsgFile     ContentType = "file"
	MsgSystem   ContentType = "system"
)

type Message struct {
	ID             uuid.UUID   `json:"id"`
	ConversationID uuid.UUID   `json:"conversation_id"`
	SenderID       uuid.UUID   `json:"sender_id"`
	Content        string      `json:"content"`
	ContentType    ContentType  `json:"content_type"`
	ReplyTo        *uuid.UUID  `json:"reply_to,omitempty"`
	DeletedAt      *time.Time  `json:"deleted_at,omitempty"`
	CreatedAt      time.Time   `json:"created_at"`
}

type SendMessageRequest struct {
	Content     string      `json:"content"`
	ContentType ContentType `json:"content_type"`
	ReplyTo     *uuid.UUID  `json:"reply_to,omitempty"`
}

type Reaction struct {
	MessageID    uuid.UUID `json:"message_id"`
	UserID       uuid.UUID `json:"user_id"`
	Emoji        string    `json:"emoji"`
	CreatedAt    time.Time `json:"created_at"`
}

type AddReactionRequest struct {
	Emoji string `json:"emoji"`
}

type Attachment struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	UploaderID     uuid.UUID `json:"uploader_id"`
	FileName       string    `json:"file_name"`
	FileSizeBytes  int64     `json:"file_size_bytes"`
	MimeType       string    `json:"mime_type"`
	StoragePath    string    `json:"storage_path"`
	CDNURL         *string   `json:"cdn_url,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}
