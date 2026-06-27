# AI Thumbnail Generation — Canvas

Canvas is Spark's AI thumbnail generation system. It creates compelling, click-optimized thumbnails that maximize viewer engagement while maintaining brand consistency and content accuracy.

## Generation Pipeline

Canvas operates in three phases:

### Phase 1 — Frame Selection
The system analyzes the full video to identify candidate frames that would make strong thumbnails:
- **Aesthetic Scoring**: A neural network trained on human-rated thumbnails scores each frame for composition, lighting, color harmony, and visual interest
- **Expression Detection**: For face-containing content, frames with high-impact facial expressions (surprise, excitement, joy, intensity) are prioritized
- **Action Detection**: Frames at motion peaks or key moments (game-winning shots, tutorial results, reaction moments) receive higher scores
- **Text Readability**: Frames with clear negative space suitable for text overlays are preferred

### Phase 2 — Enhancement
Selected frames are enhanced with generative AI:
- **Background Reframing**: Smart crop with subject-aware composition using salient object detection
- **Color Grading**: Automatic color correction and contrast enhancement optimized for thumbnail visibility
- **Super Resolution**: ESRGAN-based upscaling ensures crisp output even from low-resolution source frames
- **Text Overlay**: AI generates attention-grabbing text with optimal font, size, color contrast, and positioning

### Phase 3 — A/B Testing
Multiple thumbnail variants are generated and tested:
- Up to 5 variants per video
- Initial winner selected by engagement prediction model
- Real A/B test on first 1000 impressions settles the final choice

## Creator Control

Creators can override any AI suggestion, provide style preferences, upload custom assets, and set brand guidelines that Canvas learns and applies automatically.
