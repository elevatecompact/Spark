use async_trait::async_trait;
use serde_json::Value;

#[async_trait]
pub trait EventProducer: Send + Sync {
    async fn publish(&self, event_type: &str, data: Value);
    async fn close(&self);
}

pub struct NoopProducer;

#[async_trait]
impl EventProducer for NoopProducer {
    async fn publish(&self, event_type: &str, data: Value) {
        tracing::info!(event_type, data = %data, "noop event");
    }

    async fn close(&self) {}
}

#[cfg(feature = "kafka")]
mod kafka_impl {
    use super::*;
    use rdkafka::producer::{FutureProducer, FutureRecord};
    use rdkafka::ClientConfig;
    use std::time::Duration;

    pub struct KafkaProducer {
        producer: FutureProducer,
        topic: String,
    }

    impl KafkaProducer {
        pub fn new(brokers: &str) -> Self {
            let producer: FutureProducer = ClientConfig::new()
                .set("bootstrap.servers", brokers)
                .set("message.timeout.ms", "5000")
                .create()
                .expect("failed to create kafka producer");
            Self {
                producer,
                topic: "media-events".into(),
            }
        }
    }

    #[async_trait]
    impl EventProducer for KafkaProducer {
        async fn publish(&self, event_type: &str, data: Value) {
            let event = serde_json::json!({
                "id": uuid::Uuid::new_v4().to_string(),
                "source": "media-service",
                "specversion": "1.0",
                "type": event_type,
                "time": chrono::Utc::now().to_rfc3339(),
                "data": data,
            });
            let payload = serde_json::to_string(&event).unwrap_or_default();
            let record = FutureRecord::to(&self.topic).payload(&payload).key(&event_type);
            match self.producer.send(record, Duration::from_secs(5)).await {
                Ok(_) => tracing::info!(event_type, "event published"),
                Err((e, _)) => tracing::error!(error = %e, "failed to publish event"),
            }
        }

        async fn close(&self) {}
    }
}

#[cfg(feature = "kafka")]
pub use kafka_impl::KafkaProducer;
