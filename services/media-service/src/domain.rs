use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum ContentType {
    Video,
    Image,
    Audio,
}

impl std::fmt::Display for ContentType {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            ContentType::Video => write!(f, "video"),
            ContentType::Image => write!(f, "image"),
            ContentType::Audio => write!(f, "audio"),
        }
    }
}

impl From<&str> for ContentType {
    fn from(s: &str) -> Self {
        match s {
            "video" => ContentType::Video,
            "image" => ContentType::Image,
            "audio" => ContentType::Audio,
            _ => ContentType::Video,
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum MediaStatus {
    Uploading,
    Processing,
    Ready,
    Failed,
    Deleted,
}

impl std::fmt::Display for MediaStatus {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            MediaStatus::Uploading => write!(f, "uploading"),
            MediaStatus::Processing => write!(f, "processing"),
            MediaStatus::Ready => write!(f, "ready"),
            MediaStatus::Failed => write!(f, "failed"),
            MediaStatus::Deleted => write!(f, "deleted"),
        }
    }
}

impl From<&str> for MediaStatus {
    fn from(s: &str) -> Self {
        match s {
            "uploading" => MediaStatus::Uploading,
            "processing" => MediaStatus::Processing,
            "ready" => MediaStatus::Ready,
            "failed" => MediaStatus::Failed,
            "deleted" => MediaStatus::Deleted,
            _ => MediaStatus::Uploading,
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MediaAsset {
    pub id: Uuid,
    pub uploader_id: Uuid,
    pub content_type: ContentType,
    pub source_filename: String,
    pub file_size_bytes: i64,
    pub mime_type: String,
    pub status: MediaStatus,
    pub storage_path: String,
    pub cdn_url: String,
    pub duration_seconds: f64,
    pub width: i32,
    pub height: i32,
    pub checksum: String,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MediaRendition {
    pub id: Uuid,
    pub media_id: Uuid,
    pub profile: String,
    pub format: String,
    pub file_size_bytes: i64,
    pub storage_path: String,
    pub cdn_url: String,
    pub status: MediaStatus,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum TranscodingStatus {
    Pending,
    Processing,
    Completed,
    Failed,
}

impl std::fmt::Display for TranscodingStatus {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            TranscodingStatus::Pending => write!(f, "pending"),
            TranscodingStatus::Processing => write!(f, "processing"),
            TranscodingStatus::Completed => write!(f, "completed"),
            TranscodingStatus::Failed => write!(f, "failed"),
        }
    }
}

impl From<&str> for TranscodingStatus {
    fn from(s: &str) -> Self {
        match s {
            "pending" => TranscodingStatus::Pending,
            "processing" => TranscodingStatus::Processing,
            "completed" => TranscodingStatus::Completed,
            "failed" => TranscodingStatus::Failed,
            _ => TranscodingStatus::Pending,
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TranscodingJob {
    pub id: Uuid,
    pub media_id: Uuid,
    pub profiles: Vec<String>,
    pub status: TranscodingStatus,
    pub worker_id: String,
    pub started_at: Option<DateTime<Utc>>,
    pub completed_at: Option<DateTime<Utc>>,
    pub error_message: String,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DrmPolicy {
    pub id: Uuid,
    pub name: String,
    pub content_id: Option<Uuid>,
    pub key_system: String,
    pub license_duration_seconds: i64,
    pub security_level: String,
    pub is_active: bool,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UploadSession {
    pub id: Uuid,
    pub uploader_id: Uuid,
    pub filename: String,
    pub file_size_bytes: i64,
    pub content_type: String,
    pub chunks_total: i32,
    pub chunks_done: i32,
    pub checksum: String,
    pub status: String,
    pub storage_path: String,
    pub created_at: DateTime<Utc>,
    pub expires_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct StorageUsage {
    pub total_bytes: i64,
    pub upload_bytes: i64,
    pub rendition_bytes: i64,
    pub thumbnail_bytes: i64,
    pub asset_count: i64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct InitUploadRequest {
    pub uploader_id: Uuid,
    pub filename: String,
    pub content_type: String,
    pub file_size_bytes: i64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UploadChunkRequest {
    pub chunk_index: i32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TranscodeRequest {
    pub media_id: Uuid,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ThumbnailRequest {
    pub media_id: Uuid,
    pub time_seconds: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct OptimizeRequest {
    pub media_id: Uuid,
    pub width: i32,
    pub height: i32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DrmLicenseRequest {
    pub media_id: Uuid,
    pub key_system: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CachePurgeRequest {
    pub path: String,
}
