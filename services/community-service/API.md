# community-service — API Contract
## Communities: POST /v1/communities, GET /v1/communities/{id}, PATCH /v1/communities/{id}, DELETE /v1/communities/{id}, GET /v1/communities (discover, with filters)
## Membership: POST /v1/communities/{id}/join, POST /v1/communities/{id}/leave, GET /v1/communities/{id}/members, PATCH /v1/communities/{id}/members/{userId}/role
## Posts: POST /v1/communities/{id}/posts, GET /v1/communities/{id}/posts (paginated), PUT /v1/posts/{id}, DELETE /v1/posts/{id}, POST /v1/posts/{id}/pin, POST /v1/posts/{id}/report
## Comments: POST /v1/posts/{id}/comments, GET /v1/posts/{id}/comments, DELETE /v1/comments/{id}
## Reactions: POST /v1/posts/{id}/reactions (emoji), POST /v1/comments/{id}/reactions
## Admin: POST /v1/admin/communities/{id}/feature, POST /v1/admin/communities/{id}/suspend, GET /v1/admin/stats (platform-wide community metrics)
