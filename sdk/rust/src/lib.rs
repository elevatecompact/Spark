use reqwest::{Client as HttpClient, StatusCode};
use serde::{Deserialize, Serialize};
use thiserror::Error;

#[derive(Error, Debug)]
pub enum SparkError {
    #[error("HTTP error: {0}")]
    Http(#[from] reqwest::Error),
    #[error("API error: {status} - {message}")]
    Api { status: StatusCode, message: String },
    #[error("Authentication required")]
    Unauthorized,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct User {
    pub id: String,
    pub email: String,
    pub username: String,
    pub display_name: Option<String>,
    pub avatar_url: Option<String>,
    pub is_creator: bool,
    pub is_verified: bool,
    pub created_at: String,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct AuthResponse {
    pub user: User,
    pub access_token: String,
    pub refresh_token: String,
    pub expires_in: i64,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct Stream {
    pub id: String,
    pub creator_id: String,
    pub title: String,
    pub description: Option<String>,
    pub thumbnail_url: Option<String>,
    pub status: String,
    pub category: Option<String>,
    pub tags: Vec<String>,
    pub is_live: bool,
    pub viewer_count: i64,
    pub started_at: Option<String>,
    pub ended_at: Option<String>,
    pub created_at: String,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct Pagination {
    pub page: i32,
    pub page_size: i32,
    pub total: i32,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct StreamListResponse {
    pub streams: Vec<Stream>,
    pub pagination: Pagination,
}

pub struct Spark {
    client: HttpClient,
    base_url: String,
    access_token: Option<String>,
}

impl Spark {
    pub fn new(base_url: Option<String>, access_token: Option<String>) -> Self {
        Self {
            client: HttpClient::new(),
            base_url: base_url.unwrap_or_else(|| "https://api.spark.dev/api/v1".to_string()),
            access_token,
        }
    }

    pub fn set_access_token(&mut self, token: String) {
        self.access_token = Some(token);
    }

    async fn request<T: for<'de> Deserialize<'de>>(
        &self,
        method: reqwest::Method,
        path: &str,
        body: Option<serde_json::Value>,
    ) -> Result<T, SparkError> {
        let url = format!("{}{}", self.base_url, path);
        let mut req = self.client.request(method, url);

        if let Some(token) = &self.access_token {
            req = req.header("Authorization", format!("Bearer {}", token));
        }

        if let Some(b) = body {
            req = req.json(&b);
        }

        let resp = req.send().await?;
        let status = resp.status();

        if !status.is_success() {
            let text = resp.text().await.unwrap_or_default();
            return Err(SparkError::Api {
                status,
                message: text,
            });
        }

        Ok(resp.json::<T>().await?)
    }

    pub async fn register(
        &self,
        email: &str,
        username: &str,
        password: &str,
    ) -> Result<AuthResponse, SparkError> {
        let body = serde_json::json!({
            "email": email,
            "username": username,
            "password": password,
        });
        self.request(reqwest::Method::POST, "/auth/register", Some(body)).await
    }

    pub async fn login(&self, email: &str, password: &str) -> Result<AuthResponse, SparkError> {
        let body = serde_json::json!({
            "email": email,
            "password": password,
        });
        self.request(reqwest::Method::POST, "/auth/login", Some(body)).await
    }

    pub async fn me(&self) -> Result<User, SparkError> {
        self.request(reqwest::Method::GET, "/users/me", None).await
    }

    pub async fn get_user(&self, id: &str) -> Result<User, SparkError> {
        self.request(reqwest::Method::GET, &format!("/users/{}", id), None)
            .await
    }

    pub async fn list_streams(
        &self,
        page: Option<i32>,
        page_size: Option<i32>,
        is_live: Option<bool>,
    ) -> Result<StreamListResponse, SparkError> {
        let mut path = "/streams".to_string();
        let mut params = vec![];
        if let Some(p) = page {
            params.push(format!("page={}", p));
        }
        if let Some(ps) = page_size {
            params.push(format!("page_size={}", ps));
        }
        if let Some(l) = is_live {
            params.push(format!("is_live={}", l));
        }
        if !params.is_empty() {
            path.push('?');
            path.push_str(&params.join("&"));
        }
        self.request(reqwest::Method::GET, &path, None).await
    }

    pub async fn get_stream(&self, id: &str) -> Result<Stream, SparkError> {
        self.request(reqwest::Method::GET, &format!("/streams/{}", id), None)
            .await
    }
}
