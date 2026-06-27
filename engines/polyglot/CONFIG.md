# Configuration

| Key | Default | Description |
|-----|---------|-------------|
| model.nllb.path | /models/nllb-200-distilled | NLLB-200 model path |
| model.nllb.max_length | 256 | Max input tokens |
| model.opus.path | /models/opus-mt | Opus-MT model path |
| inference.provider | cpu | cpu, cuda, directml |
| inference.batch_size | 16 | Translation batch size |
| inference.beam_size | 4 | Beam search width |
| quality.min_confidence | 0.7 | Auto-publish threshold |
| quality.model | cometkiwi-xl | Quality estimation model |
| streaming.context_window | 5 | Prior sentences for context |
| streaming.max_latency_ms | 500 | Max streaming latency |
| cache.ttl | 3600 | Translation cache TTL |
