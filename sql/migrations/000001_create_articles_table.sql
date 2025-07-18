-- +goose Up
CREATE TABLE IF NOT EXISTS articles (
    id bigserial PRIMARY KEY,
    title text NOT NULL,
	tags text[] DEFAULT '{}'
);

-- +goose Down
DROP TABLE articles;
