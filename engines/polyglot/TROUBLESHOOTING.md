# Troubleshooting

## Translation quality poor
1. Check if custom glossary needed for domain terms.
2. Verify quality.min_confidence threshold appropriate.
3. Check active model for language pair: GET /v1/model/languages.
4. Submit quality feedback for improvement.

## Translation failing for pair
1. Verify language pair supported.
2. Check model loaded and GPU memory available.
3. For streaming, ensure source in sentence-sized chunks.
4. Check input does not exceed max_length tokens.
