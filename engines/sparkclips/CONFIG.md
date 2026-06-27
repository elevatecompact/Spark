# Configuration

| Key | Default | Description |
|-----|---------|-------------|
| modality.audio.enabled | true | Audio excitement analysis |
| modality.visual.enabled | true | Visual event detection |
| modality.chat.enabled | true | Chat sentiment analysis |
| fusion.excitement_window | 30 | Fusion window in seconds |
| fusion.weights | audio:0.3,visual:0.5,chat:0.2 | Modality weights |
| clipping.min_duration | 15 | Minimum clip duration |
| clipping.max_duration | 90 | Maximum clip duration |
| clipping.output_formats | mp4,gif | Output formats |
| clipping.padding_before | 2 | Seconds before highlight |
| clipping.padding_after | 3 | Seconds after highlight |
