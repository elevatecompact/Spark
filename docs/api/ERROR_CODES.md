# Error Codes

SPARK APIs return structured error responses with standardized error codes to help clients handle errors programmatically.

## Error Response Format

Every error response follows a consistent JSON structure. The code field contains a machine-readable error code string. The message field contains a human-readable description of the error. The details field contains additional information such as which validation rules failed. The requestId field contains a unique identifier for the failed request, useful for debugging with support.

## HTTP Status Codes

The API uses standard HTTP status codes. 200 OK indicates success. 201 Created indicates successful resource creation. 204 No Content indicates successful deletion. 400 Bad Request indicates a client error such as validation failure. 401 Unauthorized indicates missing or invalid authentication. 403 Forbidden indicates the authenticated user lacks permission. 404 Not Found indicates the requested resource does not exist. 409 Conflict indicates a conflict with the current state of the resource. 422 Unprocessable Entity indicates semantic validation errors. 429 Too Many Requests indicates rate limiting.

## Error Codes by Category

Authentication errors include INVALID_TOKEN, TOKEN_EXPIRED, INSUFFICIENT_PERMISSIONS, and ACCOUNT_SUSPENDED. Validation errors include VALIDATION_ERROR, INVALID_INPUT, MISSING_REQUIRED_FIELD, and INVALID_FORMAT. Resource errors include NOT_FOUND, ALREADY_EXISTS, CONFLICT, and RESOURCE_EXHAUSTED. Rate limiting errors include RATE_LIMIT_EXCEEDED and BURST_LIMIT_EXCEEDED. Server errors include INTERNAL_ERROR, SERVICE_UNAVAILABLE, and GATEWAY_TIMEOUT.

## Error Handling Best Practices

Clients should always check the code field rather than parsing the message string. Retry logic should use exponential backoff for 429 and 5xx errors.
