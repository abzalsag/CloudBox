CREATE TABLE files(
id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    storage_key TEXT NOT NULL UNIQUE,
    size BIGINT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updatet_at TIMESTAMP NOT NULL DEFAULT  NOW(),
);

CREATE INDEX idx_files_user_id ON files(user_id)
