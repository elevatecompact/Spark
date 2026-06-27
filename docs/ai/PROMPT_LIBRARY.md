# Prompt Library and Management

The Spark Prompt Library is a centralized repository of prompts used across all generative AI features — from content moderation explanations to thumbnail text generation and creator coaching responses.

## Library Structure

Prompts are organized by domain, model, and version:

### Domain Categories
- **Content Moderation**: Prompts for generating human-readable violation explanations, appeal responses, and policy clarification
- **Creator Tools**: Prompts for coaching recommendations, content feedback, analytics interpretations, and growth suggestions
- **Thumbnail Generation**: Prompts for text overlay generation, style suggestions, and call-to-action phrasing
- **Translation and Localization**: Prompts for cultural adaptation explanations, idiom handling, and domain-specific terminology

### Versioning and Governance

- Every prompt is versioned with a semantic version number (MAJOR.MINOR.PATCH)
- Changes require approval workflow: author to peer review to QA validation to production deploy
- Prompt changes are A/B tested against current version before full rollout
- Full audit trail maintained: who changed what, when, and why

## Prompt Engineering Standards

Prompts in the library follow standardized templates:
- **System Prompt**: Defines role, constraints, output format, and ethical guardrails
- **Context Template**: Structured placeholders for dynamic injection (e.g., video title, violation type)
- **Few-Shot Examples**: Curated examples demonstrating desired output quality and edge case handling
- **Output Schema**: JSON schema specifying structure, fields, and validations for model output

## Runtime Integration

The Prompt Library exposes a gRPC API consumed by AI services. Services request the prompt template by ID and version, inject runtime context, and send to the appropriate model. This decouples prompt management from service code — prompt updates require no code deployment.

## Testing and Evaluation

Prompt performance is continuously evaluated: output quality scores, safety compliance rate, instruction adherence, and hallucination rate. Regression tests run on every prompt change.
