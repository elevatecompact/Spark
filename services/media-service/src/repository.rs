use deadpool_postgres::Pool;
use serde_json;
use tokio_postgres::Row;
use uuid::Uuid;

use crate::domain::*;

pub struct MediaRepository {
    pool: Pool,
}

impl MediaRepository {
    pub fn new(pool: Pool) -> Self {
        Self { pool }
    }

    async fn conn(&self) -> Result<deadpool_postgres::Object, String> {
        self.pool.get().await.map_err(|e| e.to_string())
    }

    pub async fn create_media_asset(&self, m: &MediaAsset) -> Result<(), String> {
        let client = self.conn().await?;
        client
            .execute(
                "INSERT INTO media_assets (id, uploader_id, content_type, source_filename, file_size_bytes, mime_type, status, storage_path, cdn_url, duration_seconds, width, height, checksum, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)",
                &[&m.id, &m.uploader_id, &m.content_type.to_string(), &m.source_filename, &m.file_size_bytes, &m.mime_type, &m.status.to_string(), &m.storage_path, &m.cdn_url, &m.duration_seconds, &m.width, &m.height, &m.checksum, &m.created_at],
            )
            .await
            .map_err(|e| e.to_string())?;
        Ok(())
    }

    pub async fn get_media_asset(&self, id: Uuid) -> Result<MediaAsset, String> {
        let client = self.conn().await?;
        let row = client
            .query_one(
                "SELECT id, uploader_id, content_type, source_filename, file_size_bytes, mime_type, status, storage_path, cdn_url, duration_seconds, width, height, checksum, created_at FROM media_assets WHERE id=$1",
                &[&id],
            )
            .await
            .map_err(|e| e.to_string())?;
        Ok(row_to_media(&row))
    }

    pub async fn update_media_status(&self, id: Uuid, status: &MediaStatus) -> Result<(), String> {
        let client = self.conn().await?;
        client
            .execute("UPDATE media_assets SET status=$2 WHERE id=$1", &[&id, &status.to_string()])
            .await
            .map_err(|e| e.to_string())?;
        Ok(())
    }

    pub async fn update_media_cdn(&self, id: Uuid, cdn_url: &str) -> Result<(), String> {
        let client = self.conn().await?;
        client
            .execute("UPDATE media_assets SET cdn_url=$2, status='ready' WHERE id=$1", &[&id, &cdn_url])
            .await
            .map_err(|e| e.to_string())?;
        Ok(())
    }

    pub async fn delete_media(&self, id: Uuid) -> Result<(), String> {
        let client = self.conn().await?;
        client
            .execute("UPDATE media_assets SET status='deleted' WHERE id=$1", &[&id])
            .await
            .map_err(|e| e.to_string())?;
        Ok(())
    }

    pub async fn create_rendition(&self, r: &MediaRendition) -> Result<(), String> {
        let client = self.conn().await?;
        client
            .execute(
                "INSERT INTO media_renditions (id, media_id, profile, format, file_size_bytes, storage_path, cdn_url, status, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)",
                &[&r.id, &r.media_id, &r.profile, &r.format, &r.file_size_bytes, &r.storage_path, &r.cdn_url, &r.status.to_string(), &r.created_at],
            )
            .await
            .map_err(|e| e.to_string())?;
        Ok(())
    }

    pub async fn get_renditions(&self, media_id: Uuid) -> Result<Vec<MediaRendition>, String> {
        let client = self.conn().await?;
        let rows = client
            .query(
                "SELECT id, media_id, profile, format, file_size_bytes, storage_path, cdn_url, status, created_at FROM media_renditions WHERE media_id=$1 ORDER BY created_at",
                &[&media_id],
            )
            .await
            .map_err(|e| e.to_string())?;
        Ok(rows.iter().map(|row| MediaRendition {
            id: row.get("id"),
            media_id: row.get("media_id"),
            profile: row.get("profile"),
            format: row.get("format"),
            file_size_bytes: row.get("file_size_bytes"),
            storage_path: row.get("storage_path"),
            cdn_url: row.get("cdn_url"),
            status: MediaStatus::from(row.get::<_, String>("status").as_str()),
            created_at: row.get("created_at"),
        }).collect())
    }

    pub async fn create_transcoding_job(&self, j: &TranscodingJob) -> Result<(), String> {
        let client = self.conn().await?;
        let profiles = serde_json::to_value(&j.profiles).unwrap_or_default();
        client
            .execute(
                "INSERT INTO transcoding_jobs (id, media_id, profiles, status, worker_id, started_at, completed_at, error_message, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)",
                &[&j.id, &j.media_id, &profiles, &j.status.to_string(), &j.worker_id, &j.started_at, &j.completed_at, &j.error_message, &j.created_at],
            )
            .await
            .map_err(|e| e.to_string())?;
        Ok(())
    }

    pub async fn get_transcoding_job(&self, id: Uuid) -> Result<TranscodingJob, String> {
        let client = self.conn().await?;
        let row = client
            .query_one("SELECT id, media_id, profiles, status, worker_id, started_at, completed_at, error_message, created_at FROM transcoding_jobs WHERE id=$1", &[&id])
            .await
            .map_err(|e| e.to_string())?;
        Ok(row_to_job(&row))
    }

    pub async fn update_transcoding_job_status(&self, id: Uuid, status: &TranscodingStatus, err_msg: &str) -> Result<(), String> {
        let client = self.conn().await?;
        client
            .execute(
                "UPDATE transcoding_jobs SET status=$2, error_message=$3, completed_at=CASE WHEN $2 IN ('completed','failed') THEN NOW() ELSE NULL END WHERE id=$1",
                &[&id, &status.to_string(), &err_msg],
            )
            .await
            .map_err(|e| e.to_string())?;
        Ok(())
    }

    pub async fn save_upload_session(&self, s: &UploadSession) -> Result<(), String> {
        let client = self.conn().await?;
        client
            .execute(
                "INSERT INTO upload_sessions (id, uploader_id, filename, file_size_bytes, content_type, chunks_total, chunks_done, checksum, status, storage_path, created_at, expires_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)",
                &[&s.id, &s.uploader_id, &s.filename, &s.file_size_bytes, &s.content_type, &s.chunks_total, &s.chunks_done, &s.checksum, &s.status, &s.storage_path, &s.created_at, &s.expires_at],
            )
            .await
            .map_err(|e| e.to_string())?;
        Ok(())
    }

    pub async fn get_upload_session(&self, id: Uuid) -> Result<UploadSession, String> {
        let client = self.conn().await?;
        let row = client
            .query_one("SELECT id, uploader_id, filename, file_size_bytes, content_type, chunks_total, chunks_done, checksum, status, storage_path, created_at, expires_at FROM upload_sessions WHERE id=$1", &[&id])
            .await
            .map_err(|e| e.to_string())?;
        Ok(UploadSession {
            id: row.get("id"),
            uploader_id: row.get("uploader_id"),
            filename: row.get("filename"),
            file_size_bytes: row.get("file_size_bytes"),
            content_type: row.get("content_type"),
            chunks_total: row.get("chunks_total"),
            chunks_done: row.get("chunks_done"),
            checksum: row.get("checksum"),
            status: row.get("status"),
            storage_path: row.get("storage_path"),
            created_at: row.get("created_at"),
            expires_at: row.get("expires_at"),
        })
    }

    pub async fn update_upload_session(&self, id: Uuid, chunks_done: i32, status: &str) -> Result<(), String> {
        let client = self.conn().await?;
        client
            .execute("UPDATE upload_sessions SET chunks_done=$2, status=$3 WHERE id=$1", &[&id, &chunks_done, &status])
            .await
            .map_err(|e| e.to_string())?;
        Ok(())
    }

    pub async fn delete_upload_session(&self, id: Uuid) -> Result<(), String> {
        let client = self.conn().await?;
        client.execute("DELETE FROM upload_sessions WHERE id=$1", &[&id]).await.map_err(|e| e.to_string())?;
        Ok(())
    }

    pub async fn create_drm_policy(&self, p: &DrmPolicy) -> Result<(), String> {
        let client = self.conn().await?;
        client
            .execute(
                "INSERT INTO drm_policies (id, name, content_id, key_system, license_duration_seconds, security_level, is_active, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
                &[&p.id, &p.name, &p.content_id, &p.key_system, &p.license_duration_seconds, &p.security_level, &p.is_active, &p.created_at],
            )
            .await
            .map_err(|e| e.to_string())?;
        Ok(())
    }

    pub async fn get_drm_policies(&self) -> Result<Vec<DrmPolicy>, String> {
        let client = self.conn().await?;
        let rows = client
            .query("SELECT id, name, content_id, key_system, license_duration_seconds, security_level, is_active, created_at FROM drm_policies WHERE is_active=true", &[])
            .await
            .map_err(|e| e.to_string())?;
        Ok(rows.iter().map(|row| DrmPolicy {
            id: row.get("id"),
            name: row.get("name"),
            content_id: row.get("content_id"),
            key_system: row.get("key_system"),
            license_duration_seconds: row.get("license_duration_seconds"),
            security_level: row.get("security_level"),
            is_active: row.get("is_active"),
            created_at: row.get("created_at"),
        }).collect())
    }

    pub async fn get_storage_usage(&self) -> Result<StorageUsage, String> {
        let client = self.conn().await?;
        let total: i64 = client.query_one("SELECT COALESCE(SUM(file_size_bytes),0) FROM media_assets WHERE status!='deleted'", &[]).await.map_err(|e| e.to_string())?.get(0);
        let count: i64 = client.query_one("SELECT COUNT(*) FROM media_assets WHERE status!='deleted'", &[]).await.map_err(|e| e.to_string())?.get(0);
        Ok(StorageUsage {
            total_bytes: total,
            upload_bytes: total,
            rendition_bytes: 0,
            thumbnail_bytes: 0,
            asset_count: count,
        })
    }

    pub async fn list_transcoding_jobs(&self, status: Option<&str>, limit: i64, offset: i64) -> Result<Vec<TranscodingJob>, String> {
        let client = self.conn().await?;
        let rows = if let Some(s) = status {
            client
                .query("SELECT id, media_id, profiles, status, worker_id, started_at, completed_at, error_message, created_at FROM transcoding_jobs WHERE status=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3", &[&s, &limit, &offset])
                .await
                .map_err(|e| e.to_string())?
        } else {
            client
                .query("SELECT id, media_id, profiles, status, worker_id, started_at, completed_at, error_message, created_at FROM transcoding_jobs ORDER BY created_at DESC LIMIT $1 OFFSET $2", &[&limit, &offset])
                .await
                .map_err(|e| e.to_string())?
        };
        Ok(rows.iter().map(|row| row_to_job(row)).collect())
    }
}

fn row_to_media(row: &Row) -> MediaAsset {
    MediaAsset {
        id: row.get("id"),
        uploader_id: row.get("uploader_id"),
        content_type: ContentType::from(row.get::<_, String>("content_type").as_str()),
        source_filename: row.get("source_filename"),
        file_size_bytes: row.get("file_size_bytes"),
        mime_type: row.get("mime_type"),
        status: MediaStatus::from(row.get::<_, String>("status").as_str()),
        storage_path: row.get("storage_path"),
        cdn_url: row.get("cdn_url"),
        duration_seconds: row.get("duration_seconds"),
        width: row.get("width"),
        height: row.get("height"),
        checksum: row.get("checksum"),
        created_at: row.get("created_at"),
    }
}

fn row_to_job(row: &Row) -> TranscodingJob {
    let profiles_val: serde_json::Value = row.get("profiles");
    let profiles: Vec<String> = serde_json::from_value(profiles_val).unwrap_or_default();
    TranscodingJob {
        id: row.get("id"),
        media_id: row.get("media_id"),
        profiles,
        status: TranscodingStatus::from(row.get::<_, String>("status").as_str()),
        worker_id: row.get("worker_id"),
        started_at: row.get("started_at"),
        completed_at: row.get("completed_at"),
        error_message: row.get("error_message"),
        created_at: row.get("created_at"),
    }
}
