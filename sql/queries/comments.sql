-- name: CreateComment :many
WITH ins AS (
    INSERT INTO comments (id, article_id, username, content)
    VALUES ($1, $2, $3, $4)
    RETURNING article_id
)
SELECT c.id, c.article_id, c.username, c.content
FROM comments AS c
JOIN ins ON ins.article_id = c.article_id
ORDER BY c.id DESC;
