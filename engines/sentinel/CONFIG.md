# Configuration

| Key | Default | Description |
|-----|---------|-------------|
| tier1.blocklist.path | /data/blocklists/ | Path to blocklist files |
| tier2.text.model | distilbert-base-uncased | Lightweight text classifier |
| tier2.image.model | resnet-18 | Lightweight image classifier |
| tier2.confidence_threshold | 0.85 | Auto-block threshold |
| tier3.text.model | roberta-large | Deep text analysis model |
| tier3.image.model | clip-vit-large | Multimodal analysis model |
| tier3.budget_per_request | 0.001 | Max USD budget for tier 3 |
| video.sample_interval | 30 | Sample every Nth frame |
| active_learning.sample_rate | 0.01 | Fraction for human review |
