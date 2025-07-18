-- name: GetArticleByID :one
SELECT id, title, tags
FROM articles
WHERE id = $1;
