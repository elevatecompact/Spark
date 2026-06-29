pub struct Config {
    pub port: String,
    pub database_url: String,
    pub kafka_brokers: String,
    pub log_level: String,
}

impl Config {
    pub fn load() -> Self {
        Self {
            port: std::env::var("MEDIA_PORT").unwrap_or_else(|_| "4024".into()),
            database_url: std::env::var("MEDIA_DB_URL").unwrap_or_else(|_| "postgres://spark:spark@localhost:5432/spark_media?sslmode=disable".into()),
            kafka_brokers: std::env::var("MEDIA_KAFKA_BROKERS").unwrap_or_else(|_| "localhost:9092".into()),
            log_level: std::env::var("LOG_LEVEL").unwrap_or_else(|_| "info".into()),
        }
    }
}
