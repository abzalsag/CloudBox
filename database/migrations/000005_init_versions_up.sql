CREATE TABLE file_versions (
                               id             UUID PRIMARY KEY,
                               file_id        UUID        NOT NULL REFERENCES files (id) ON DELETE CASCADE,
                               version_number INT         NOT NULL,
                               size           BIGINT      NOT NULL,
                               storage_key    TEXT        NOT NULL,
                               checksum       TEXT        NOT NULL,
                               created_by     UUID        NOT NULL REFERENCES users (id),
                               created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
                               UNIQUE (file_id, version_number)
);