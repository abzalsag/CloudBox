CREATE TABLE files (
                       id          UUID PRIMARY KEY,
                       owner_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
                       folder_id   UUID        REFERENCES folders (id) ON DELETE CASCADE,
                       name        TEXT        NOT NULL,
                       size        BIGINT      NOT NULL,
                       mime_type   TEXT        NOT NULL DEFAULT 'application/octet-stream',
                       checksum    TEXT        NOT NULL,
                       storage_key TEXT        NOT NULL,
                       bucket      TEXT        NOT NULL,
                       created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
                       updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
                       deleted_at  TIMESTAMPTZ,
                       CONSTRAINT files_name_not_empty CHECK (length(btrim(name)) > 0),
                       CONSTRAINT files_size_nonneg CHECK (size >= 0)
    );

CREATE UNIQUE INDEX files_unique_name_in_folder
    ON files (owner_id, folder_id, name)
    WHERE deleted_at IS NULL AND folder_id IS NOT NULL;

CREATE UNIQUE INDEX files_unique_name_in_root
    ON files (owner_id, name)
    WHERE deleted_at IS NULL AND folder_id IS NULL;

CREATE INDEX files_owner_folder_idx ON files (owner_id, folder_id);
CREATE INDEX files_checksum_idx ON files (checksum);
CREATE INDEX files_deleted_at_idx ON files (deleted_at) WHERE deleted_at IS NOT NULL;
CREATE UNIQUE INDEX files_bucket_storage_key_key ON files (bucket, storage_key);