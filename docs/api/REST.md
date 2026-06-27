# REST API

The SPARK REST API provides a straightforward HTTP-based interface for third-party integrations, simple CRUD operations, and scenarios where GraphQL is not available or practical.

## Base URL

All REST API endpoints are available at https://api.sparkplatform.com/v1/. Every response includes a Content-Type header of application/json. Requests must include the Content-Type header for methods with a request body.

## Resources

Resources follow a hierarchical naming convention. Top-level resources include /users, /content, /subscriptions, /transactions, and /notifications. Nested resources include /users/{userId}/content, /content/{contentId}/comments, and /users/{userId}/subscriptions.

## HTTP Methods

GET requests are idempotent and safe. POST requests create new resources and return 201 Created with the resource location in the Location header. PATCH requests perform partial updates and return 200 OK with the updated resource. DELETE requests remove resources and return 204 No Content. PUT is not used; all updates use PATCH for partial updates.

## Response Format

Successful responses return the requested resource or collection. Error responses follow the standard error format with code, message, and details fields. Paginated responses include a meta object with cursor, hasMore, and totalCount fields.

## Conditional Requests

ETags support conditional GET requests with If-None-Match headers for caching. If-Modified-Since headers are also supported. Immutable resources may be cached for up to 24 hours with appropriate Cache-Control headers.
