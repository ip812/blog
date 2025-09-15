-- name: CreateComment :one
INSERT INTO comments (id, article_id, username, content)
VALUES ($1, $2, $3, $4)
RETURNING id, article_id, username, content;

-- name: GetCommentsByArticleId :many
SELECT id, article_id, username, content
FROM comments
WHERE article_id = $1;
