use std::sync::Arc;
use chrono::{Duration, Utc};
use rand::rngs::StdRng;
use rand::{Rng, SeedableRng};
use tokio::time::sleep;
use uuid::Uuid;

use crate::domain::*;
use crate::events::EventProducer;
use crate::repository::MediaRepository;

pub struct MediaService {
    repo: Arc<MediaRepository>,
    pub evt: Arc<dyn EventProducer + Send + Sync>,
}

impl MediaService {
    pub fn new(repo: Arc<MediaRepository>, evt: Arc<dyn EventProducer + Send + Sync>) -> Self {
        Self { repo, evt }
    }

    pub async fn init_upload(&self, req: &InitUploadRequest) -> Result<UploadSession, String> {
        let chunk_size: i64 = 5 * 1024 * 1024;
        let chunks_total = ((req.file_size_bytes + chunk_size - 1) / chunk_size) as i32;
        let session = UploadSession {
            id: Uuid::new_v4(),
            uploader_id: req.uploader_id,
            filename: req.filename.clone(),
            file_size_bytes: req.file_size_bytes,
            content_type: req.content_type.clone(),
            chunks_total,
            chunks_done: 0,
            checksum: String::new(),
            status: "initiated".into(),
            storage_path: format!("uploads/{}/{}", req.uploader_id, Uuid::new_v4()),
            created_at: Utc::now(),
            expires_at: Utc::now() + Duration::hours(24),
        };
        self.repo.save_upload_session(&session).await.map_err(|e| e.to_string())?;
        Ok(session)
    }

    pub async fn upload_chunk(&self, session_id: Uuid) -> Result<(), String> {
        let session = self.repo.get_upload_session(session_id).await.map_err(|e| e.to_string())?;
        if session.status == "completed" {
            return Err("upload already completed".into());
        }
        let done = session.chunks_done + 1;
        let status = if done >= session.chunks_total { "completed" } else { "uploading" };
        self.repo.update_upload_session(session_id, done, status).await.map_err(|e| e.to_string())?;
        Ok(())
    }

    pub async fn complete_upload(&self, session_id: Uuid) -> Result<MediaAsset, String> {
        let session = self.repo.get_upload_session(session_id).await.map_err(|e| e.to_string())?;
        if session.chunks_done < session.chunks_total {
            return Err("upload not complete".into());
        }
        let asset = MediaAsset {
            id: Uuid::new_v4(),
            uploader_id: session.uploader_id,
            content_type: classify_content_type(&session.content_type),
            source_filename: session.filename.clone(),
            file_size_bytes: session.file_size_bytes,
            mime_type: session.content_type.clone(),
            status: MediaStatus::Processing,
            storage_path: session.storage_path.clone(),
            cdn_url: format!("https://cdn.spark.dev/media/{}", Uuid::new_v4()),
            duration_seconds: 0.0,
            width: 1920,
            height: 1080,
            checksum: session.checksum.clone(),
            created_at: Utc::now(),
        };
        self.repo.create_media_asset(&asset).await.map_err(|e| e.to_string())?;
        self.repo.delete_upload_session(session_id).await.ok();
        let _ = self.start_transcoding(asset.id).await;
        self.evt.publish("media.upload.completed", serde_json::json!({
            "mediaId": asset.id, "uploaderId": asset.uploader_id, "fileSize": asset.file_size_bytes
        })).await;
        Ok(asset)
    }

    pub async fn get_upload_status(&self, session_id: Uuid) -> Result<UploadSession, String> {
        self.repo.get_upload_session(session_id).await.map_err(|e| e.to_string())
    }

    pub async fn cancel_upload(&self, session_id: Uuid) -> Result<(), String> {
        self.repo.delete_upload_session(session_id).await.map_err(|e| e.to_string())
    }

    pub async fn start_transcoding(&self, media_id: Uuid) -> Result<(), String> {
        let media = self.repo.get_media_asset(media_id).await.map_err(|e| e.to_string())?;
        if media.content_type != ContentType::Video {
            return Ok(());
        }
        let profiles = vec!["720p".into(), "1080p".into(), "source".into()];
        let job = TranscodingJob {
            id: Uuid::new_v4(),
            media_id,
            profiles: profiles.clone(),
            status: TranscodingStatus::Pending,
            worker_id: String::new(),
            started_at: None,
            completed_at: None,
            error_message: String::new(),
            created_at: Utc::now(),
        };
        self.repo.create_transcoding_job(&job).await.map_err(|e| e.to_string())?;
        self.evt.publish("media.transcoding.started", serde_json::json!({
            "mediaId": media_id, "jobId": job.id, "profiles": profiles
        })).await;
        let repo = self.repo.clone();
        let evt = self.evt.clone();
        tokio::spawn(async move {
            process_transcoding(repo, evt, job.id).await;
        });
        Ok(())
    }

    pub async fn get_media_status(&self, media_id: Uuid) -> Result<MediaAsset, String> {
        self.repo.get_media_asset(media_id).await.map_err(|e| e.to_string())
    }

    pub async fn get_media_info(&self, media_id: Uuid) -> Result<MediaAsset, String> {
        self.repo.get_media_asset(media_id).await.map_err(|e| e.to_string())
    }

    pub async fn get_playback_url(&self, media_id: Uuid) -> Result<String, String> {
        let media = self.repo.get_media_asset(media_id).await.map_err(|e| e.to_string())?;
        if media.status != MediaStatus::Ready {
            return Err("media not ready".into());
        }
        Ok(format!("https://cdn.spark.dev/media/{}/playlist.m3u8", media_id))
    }

    pub async fn get_download_url(&self, media_id: Uuid) -> Result<String, String> {
        Ok(format!("https://cdn.spark.dev/media/{}/download", media_id))
    }

    pub async fn get_thumbnail_url(&self, media_id: Uuid, time_secs: f64) -> Result<String, String> {
        Ok(format!("https://cdn.spark.dev/media/{}/thumb/{}.jpg", media_id, time_secs))
    }

    pub async fn generate_thumbnail(&self, req: &ThumbnailRequest) -> Result<MediaRendition, String> {
        let mut rng = StdRng::from_entropy();
        let rendition = MediaRendition {
            id: Uuid::new_v4(),
            media_id: req.media_id,
            profile: "thumbnail".into(),
            format: "jpg".into(),
            file_size_bytes: rng.gen_range(10000..50000),
            storage_path: format!("thumbnails/{}/{}.jpg", req.media_id, req.time_seconds),
            cdn_url: format!("https://cdn.spark.dev/media/{}/thumb/{}.jpg", req.media_id, req.time_seconds),
            status: MediaStatus::Ready,
            created_at: Utc::now(),
        };
        self.repo.create_rendition(&rendition).await.map_err(|e| e.to_string())?;
        self.evt.publish("media.thumbnail.generated", serde_json::json!({
            "mediaId": req.media_id, "time": req.time_seconds, "url": rendition.cdn_url
        })).await;
        Ok(rendition)
    }

    pub async fn optimize_image(&self, req: &OptimizeRequest) -> Result<MediaRendition, String> {
        let media = self.repo.get_media_asset(req.media_id).await.map_err(|e| e.to_string())?;
        if media.content_type != ContentType::Image {
            return Err("optimization only supported for images".into());
        }
        let mut rng = StdRng::from_entropy();
        let rendition = MediaRendition {
            id: Uuid::new_v4(),
            media_id: req.media_id,
            profile: format!("{}x{}", req.width, req.height),
            format: "webp".into(),
            file_size_bytes: rng.gen_range(10000..100000),
            storage_path: format!("optimized/{}/{}x{}.webp", req.media_id, req.width, req.height),
            cdn_url: format!("https://cdn.spark.dev/media/{}/opt/{}x{}.webp", req.media_id, req.width, req.height),
            status: MediaStatus::Ready,
            created_at: Utc::now(),
        };
        self.repo.create_rendition(&rendition).await.map_err(|e| e.to_string())?;
        Ok(rendition)
    }

    pub async fn get_renditions(&self, media_id: Uuid) -> Result<Vec<MediaRendition>, String> {
        self.repo.get_renditions(media_id).await.map_err(|e| e.to_string())
    }

    pub async fn issue_drm_license(&self, req: &DrmLicenseRequest) -> Result<String, String> {
        self.evt.publish("media.drm.license.issued", serde_json::json!({
            "mediaId": req.media_id, "keySystem": req.key_system
        })).await;
        Ok(format!("drm-license-{}-{}", req.media_id, Uuid::new_v4()))
    }

    pub async fn create_drm_policy(&self, p: &DrmPolicy) -> Result<DrmPolicy, String> {
        let mut policy = p.clone();
        policy.id = Uuid::new_v4();
        policy.created_at = Utc::now();
        self.repo.create_drm_policy(&policy).await.map_err(|e| e.to_string())?;
        Ok(policy)
    }

    pub async fn get_drm_policies(&self) -> Result<Vec<DrmPolicy>, String> {
        self.repo.get_drm_policies().await.map_err(|e| e.to_string())
    }

    pub async fn get_storage_usage(&self) -> Result<StorageUsage, String> {
        self.repo.get_storage_usage().await.map_err(|e| e.to_string())
    }

    pub async fn purge_cdn(&self, path: &str) -> Result<(), String> {
        tracing::info!(path, "CDN cache purged (noop)");
        Ok(())
    }

    pub async fn get_processing_queue(&self, limit: i64, offset: i64) -> Result<Vec<TranscodingJob>, String> {
        self.repo.list_transcoding_jobs(Some("pending"), limit, offset).await.map_err(|e| e.to_string())
    }

    pub async fn retry_transcoding(&self, media_id: Uuid) -> Result<(), String> {
        self.start_transcoding(media_id).await
    }
}

async fn process_transcoding(repo: Arc<MediaRepository>, evt: Arc<dyn EventProducer + Send + Sync>, job_id: Uuid) {
    let _ = repo.update_transcoding_job_status(job_id, &TranscodingStatus::Processing, "").await;
    sleep(std::time::Duration::from_millis(500)).await;
    let job = repo.get_transcoding_job(job_id).await;
    if let Ok(job) = job {
    let mut rng = StdRng::from_entropy();
    for profile in &job.profiles {
            let format = if profile == "thumbnail" { "jpg" } else { "hls" };
            let rendition = MediaRendition {
                id: Uuid::new_v4(),
                media_id: job.media_id,
                profile: profile.clone(),
                format: format.into(),
                file_size_bytes: rng.gen_range(1000000..50000000),
                storage_path: format!("renditions/{}/{}", job.media_id, profile),
                cdn_url: format!("https://cdn.spark.dev/media/{}/{}", job.media_id, profile),
                status: MediaStatus::Ready,
                created_at: Utc::now(),
            };
            let _ = repo.create_rendition(&rendition).await;
        }
        let _ = repo.update_transcoding_job_status(job_id, &TranscodingStatus::Completed, "").await;
        let _ = repo.update_media_cdn(job.media_id, &format!("https://cdn.spark.dev/media/{}/playlist.m3u8", job.media_id)).await;
        evt.publish("media.transcoding.completed", serde_json::json!({
            "mediaId": job.media_id, "jobId": job.id, "profiles": job.profiles
        })).await;
    }
}

fn classify_content_type(mime: &str) -> ContentType {
    match mime {
        "video/mp4" | "video/webm" | "video/quicktime" => ContentType::Video,
        "image/jpeg" | "image/png" | "image/webp" | "image/gif" => ContentType::Image,
        "audio/mpeg" | "audio/ogg" | "audio/wav" => ContentType::Audio,
        _ => ContentType::Video,
    }
}
