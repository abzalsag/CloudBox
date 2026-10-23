CREATE TABLE users (
                       id            UUID PRIMARY KEY,
                       email         TEXT        NOT NULL,
                       password_hash TEXT        NOT NULL,
                       full_name     TEXT        NOT NULL,
                       storage_quota BIGINT      NOT NULL DEFAULT 5368709120,  
                       storage_used  BIGINT      NOT NULL DEFAULT 0,
                       created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
                       updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
                       CONSTRAINT users_storage_used_nonneg CHECK (storage_used >= 0),
                       CONSTRAINT users_storage_quota_nonneg CHECK (storage_quota >= 0)
);

CREATE UNIQUE INDEX users_email_key ON users (lower(email));