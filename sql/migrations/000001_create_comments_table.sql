-- +goose Up
CREATE TABLE IF NOT EXISTS comments (
    id bigserial PRIMARY KEY,
    article_id bigserial NOT NULL,
    username text NOT NULL,
    content text NOT NULL
);

-- +goose Down
DROP TABLE comments;
