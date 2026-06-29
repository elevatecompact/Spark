mod config;
mod domain;
mod events;
mod repository;
mod service;

use std::net::SocketAddr;
use std::sync::Arc;

use axum::extract::{Path, Query, State};
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::routing::{delete, get, post};
use axum::{Json, Router};
use serde::Deserialize;
use tower_http::cors::CorsLayer;
use tower_http::trace::TraceLayer;
use uuid::Uuid;

use config::Config;
use domain::*;
use events::{EventProducer, NoopProducer};
use repository::MediaRepository;
use service::MediaService;

#[cfg(feature = "kafka")]
use events::KafkaProducer;

struct AppState {
    svc: MediaService,
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let cfg = Config::load();

    tracing_subscriber::fmt()
        .with_env_filter(&cfg.log_level)
        .init();

    let mut deadpool_cfg = deadpool_postgres::Config::new();
    deadpool_cfg.url = Some(cfg.database_url.clone());
    let pool = deadpool_cfg
        .create_pool(Some(deadpool_postgres::Runtime::Tokio1), tokio_postgres::NoTls)
        .expect("failed to create pool");

    let repo = Arc::new(MediaRepository::new(pool));

    let evt: Arc<dyn EventProducer + Send + Sync> = {
        #[cfg(feature = "kafka")]
        {
            if cfg.kafka_brokers.is_empty() {
                Arc::new(NoopProducer)
            } else {
                Arc::new(KafkaProducer::new(&cfg.kafka_brokers))
            }
        }
        #[cfg(not(feature = "kafka"))]
        Arc::new(NoopProducer)
    };

    let svc = MediaService::new(repo, evt);

    let state = Arc::new(AppState { svc });

    let app = Router::new()
        // Upload
        .route("/v1/upload/init", post(init_upload))
        .route("/v1/upload/{id}/chunk", post(upload_chunk))
        .route("/v1/upload/{id}/complete", post(complete_upload))
        .route("/v1/upload/{id}/status", get(upload_status))
        .route("/v1/upload/{id}", delete(cancel_upload))
        // Media
        .route("/v1/media/transcode", post(start_transcoding))
        .route("/v1/media/thumbnail", post(generate_thumbnail))
        .route("/v1/media/optimize", post(optimize_image))
        .route("/v1/media/{id}/status", get(media_status))
        .route("/v1/media/{id}/playback", get(get_playback))
        .route("/v1/media/{id}/thumbnail/{time}", get(get_thumbnail))
        .route("/v1/media/{id}/download", get(get_download))
        .route("/v1/media/{id}/info", get(media_info))
        .route("/v1/media/{id}/renditions", get(get_renditions))
        // DRM
        .route("/v1/drm/license", post(issue_license))
        .route("/v1/drm/policy", post(create_drm_policy))
        .route("/v1/drm/policies", get(get_drm_policies))
        // Admin
        .route("/v1/admin/storage/usage", get(storage_usage))
        .route("/v1/admin/cache/purge", post(purge_cdn))
        .route("/v1/admin/processing/queue", get(processing_queue))
        .route("/v1/admin/media/{id}/retry", post(retry_transcoding))
        .layer(CorsLayer::permissive())
        .layer(TraceLayer::new_for_http())
        .with_state(state);

    let addr: SocketAddr = format!("0.0.0.0:{}", cfg.port).parse()?;
    tracing::info!("media-service listening on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await?;
    axum::serve(listener, app).await?;

    Ok(())
}

// --- Upload handlers ---

async fn init_upload(
    State(state): State<Arc<AppState>>,
    Json(body): Json<InitUploadRequest>,
) -> impl IntoResponse {
    match state.svc.init_upload(&body).await {
        Ok(session) => (StatusCode::CREATED, Json(serde_json::to_value(session).unwrap())).into_response(),
        Err(e) => (StatusCode::INTERNAL_SERVER_ERROR, Json(serde_json::json!({"error": e}))).into_response(),
    }
}

async fn upload_chunk(
    State(state): State<Arc<AppState>>,
    Path(id): Path<Uuid>,
    Json(_body): Json<UploadChunkRequest>,
) -> impl IntoResponse {
    match state.svc.upload_chunk(id).await {
        Ok(_) => (StatusCode::OK, Json(serde_json::json!({"status": "chunk received"}))).into_response(),
        Err(e) => (StatusCode::INTERNAL_SERVER_ERROR, Json(serde_json::json!({"error": e}))).into_response(),
    }
}

async fn complete_upload(
    State(state): State<Arc<AppState>>,
    Path(id): Path<Uuid>,
) -> impl IntoResponse {
    match state.svc.complete_upload(id).await {
        Ok(asset) => (StatusCode::OK, Json(serde_json::to_value(asset).unwrap())).into_response(),
        Err(e) => (StatusCode::INTERNAL_SERVER_ERROR, Json(serde_json::json!({"error": e}))).into_response(),
    }
}

async fn upload_status(
    State(state): State<Arc<AppState>>,
    Path(id): Path<Uuid>,
) -> impl IntoResponse {
    match state.svc.get_upload_status(id).await {
        Ok(session) => (StatusCode::OK, Json(serde_json::to_value(session).unwrap())).into_response(),
        Err(_) => (StatusCode::NOT_FOUND, Json(serde_json::json!({"error": "upload not found"}))).into_response(),
    }
}

async fn cancel_upload(
    State(state): State<Arc<AppState>>,
    Path(id): Path<Uuid>,
) -> impl IntoResponse {
    match state.svc.cancel_upload(id).await {
        Ok(_) => (StatusCode::OK, Json(serde_json::json!({"status": "cancelled"}))).into_response(),
        Err(e) => (StatusCode::INTERNAL_SERVER_ERROR, Json(serde_json::json!({"error": e}))).into_response(),
    }
}

// --- Media handlers ---

async fn start_transcoding(
    State(state): State<Arc<AppState>>,
    Json(body): Json<TranscodeRequest>,
) -> impl IntoResponse {
    match state.svc.start_transcoding(body.media_id).await {
        Ok(_) => (StatusCode::ACCEPTED, Json(serde_json::json!({"status": "transcoding started"}))).into_response(),
        Err(e) => (StatusCode::INTERNAL_SERVER_ERROR, Json(serde_json::json!({"error": e}))).into_response(),
    }
}

async fn generate_thumbnail(
    State(state): State<Arc<AppState>>,
    Json(body): Json<ThumbnailRequest>,
) -> impl IntoResponse {
    match state.svc.generate_thumbnail(&body).await {
        Ok(rd) => (StatusCode::CREATED, Json(serde_json::to_value(rd).unwrap())).into_response(),
        Err(e) => (StatusCode::INTERNAL_SERVER_ERROR, Json(serde_json::json!({"error": e}))).into_response(),
    }
}

async fn optimize_image(
    State(state): State<Arc<AppState>>,
    Json(body): Json<OptimizeRequest>,
) -> impl IntoResponse {
    match state.svc.optimize_image(&body).await {
        Ok(rd) => (StatusCode::CREATED, Json(serde_json::to_value(rd).unwrap())).into_response(),
        Err(e) => (StatusCode::INTERNAL_SERVER_ERROR, Json(serde_json::json!({"error": e}))).into_response(),
    }
}

async fn media_status(
    State(state): State<Arc<AppState>>,
    Path(id): Path<Uuid>,
) -> impl IntoResponse {
    match state.svc.get_media_status(id).await {
        Ok(media) => (StatusCode::OK, Json(serde_json::to_value(media).unwrap())).into_response(),
        Err(_) => (StatusCode::NOT_FOUND, Json(serde_json::json!({"error": "media not found"}))).into_response(),
    }
}

async fn get_playback(
    State(state): State<Arc<AppState>>,
    Path(id): Path<Uuid>,
) -> impl IntoResponse {
    match state.svc.get_playback_url(id).await {
        Ok(url) => (StatusCode::OK, Json(serde_json::json!({"url": url}))).into_response(),
        Err(e) => (StatusCode::NOT_FOUND, Json(serde_json::json!({"error": e}))).into_response(),
    }
}

async fn get_thumbnail(
    State(state): State<Arc<AppState>>,
    Path((id, time)): Path<(Uuid, f64)>,
) -> impl IntoResponse {
    match state.svc.get_thumbnail_url(id, time).await {
        Ok(url) => (StatusCode::OK, Json(serde_json::json!({"url": url}))).into_response(),
        Err(e) => (StatusCode::INTERNAL_SERVER_ERROR, Json(serde_json::json!({"error": e}))).into_response(),
    }
}

async fn get_download(
    State(state): State<Arc<AppState>>,
    Path(id): Path<Uuid>,
) -> impl IntoResponse {
    match state.svc.get_download_url(id).await {
        Ok(url) => (StatusCode::OK, Json(serde_json::json!({"url": url}))).into_response(),
        Err(e) => (StatusCode::NOT_FOUND, Json(serde_json::json!({"error": e}))).into_response(),
    }
}

async fn media_info(
    State(state): State<Arc<AppState>>,
    Path(id): Path<Uuid>,
) -> impl IntoResponse {
    match state.svc.get_media_info(id).await {
        Ok(media) => (StatusCode::OK, Json(serde_json::to_value(media).unwrap())).into_response(),
        Err(_) => (StatusCode::NOT_FOUND, Json(serde_json::json!({"error": "media not found"}))).into_response(),
    }
}

async fn get_renditions(
    State(state): State<Arc<AppState>>,
    Path(id): Path<Uuid>,
) -> impl IntoResponse {
    match state.svc.get_renditions(id).await {
        Ok(renditions) => (StatusCode::OK, Json(serde_json::to_value(renditions).unwrap())).into_response(),
        Err(e) => (StatusCode::INTERNAL_SERVER_ERROR, Json(serde_json::json!({"error": e}))).into_response(),
    }
}

// --- DRM handlers ---

async fn issue_license(
    State(state): State<Arc<AppState>>,
    Json(body): Json<DrmLicenseRequest>,
) -> impl IntoResponse {
    match state.svc.issue_drm_license(&body).await {
        Ok(license) => (StatusCode::OK, Json(serde_json::json!({"license": license}))).into_response(),
        Err(e) => (StatusCode::INTERNAL_SERVER_ERROR, Json(serde_json::json!({"error": e}))).into_response(),
    }
}

async fn create_drm_policy(
    State(state): State<Arc<AppState>>,
    Json(body): Json<DrmPolicy>,
) -> impl IntoResponse {
    match state.svc.create_drm_policy(&body).await {
        Ok(policy) => (StatusCode::CREATED, Json(serde_json::to_value(policy).unwrap())).into_response(),
        Err(e) => (StatusCode::INTERNAL_SERVER_ERROR, Json(serde_json::json!({"error": e}))).into_response(),
    }
}

async fn get_drm_policies(
    State(state): State<Arc<AppState>>,
) -> impl IntoResponse {
    match state.svc.get_drm_policies().await {
        Ok(policies) => (StatusCode::OK, Json(serde_json::to_value(policies).unwrap())).into_response(),
        Err(e) => (StatusCode::INTERNAL_SERVER_ERROR, Json(serde_json::json!({"error": e}))).into_response(),
    }
}

// --- Admin handlers ---

async fn storage_usage(
    State(state): State<Arc<AppState>>,
) -> impl IntoResponse {
    match state.svc.get_storage_usage().await {
        Ok(usage) => (StatusCode::OK, Json(serde_json::to_value(usage).unwrap())).into_response(),
        Err(e) => (StatusCode::INTERNAL_SERVER_ERROR, Json(serde_json::json!({"error": e}))).into_response(),
    }
}

#[derive(Deserialize)]
struct PurgeBody {
    path: String,
}

async fn purge_cdn(
    State(state): State<Arc<AppState>>,
    Json(body): Json<PurgeBody>,
) -> impl IntoResponse {
    match state.svc.purge_cdn(&body.path).await {
        Ok(_) => (StatusCode::OK, Json(serde_json::json!({"status": "purged"}))).into_response(),
        Err(e) => (StatusCode::INTERNAL_SERVER_ERROR, Json(serde_json::json!({"error": e}))).into_response(),
    }
}

#[derive(Deserialize)]
struct QueueParams {
    limit: Option<i64>,
    offset: Option<i64>,
}

async fn processing_queue(
    State(state): State<Arc<AppState>>,
    Query(params): Query<QueueParams>,
) -> impl IntoResponse {
    let limit = params.limit.unwrap_or(20).min(100).max(1);
    let offset = params.offset.unwrap_or(0).max(0);
    match state.svc.get_processing_queue(limit, offset).await {
        Ok(jobs) => (StatusCode::OK, Json(serde_json::to_value(jobs).unwrap())).into_response(),
        Err(e) => (StatusCode::INTERNAL_SERVER_ERROR, Json(serde_json::json!({"error": e}))).into_response(),
    }
}

async fn retry_transcoding(
    State(state): State<Arc<AppState>>,
    Path(id): Path<Uuid>,
) -> impl IntoResponse {
    match state.svc.retry_transcoding(id).await {
        Ok(_) => (StatusCode::OK, Json(serde_json::json!({"status": "retrying"}))).into_response(),
        Err(e) => (StatusCode::INTERNAL_SERVER_ERROR, Json(serde_json::json!({"error": e}))).into_response(),
    }
}
