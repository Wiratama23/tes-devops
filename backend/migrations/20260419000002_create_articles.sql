-- +goose Up
CREATE TABLE articles (
    articles_id SERIAL PRIMARY KEY,
    uid UUID NOT NULL REFERENCES users(uid) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    article_text TEXT NOT NULL,
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_articles_uid ON articles(uid);

-- +goose Down
DROP TABLE IF EXISTS articles;
