CREATE TYPE resource_type AS ENUM ('file', 'folder');
CREATE TYPE share_permission AS ENUM ('viewer', 'editor');

CREATE TABLE shares (
                        id                  UUID PRIMARY KEY,
                        resource_type       resource_type    NOT NULL,
                        resource_id         UUID             NOT NULL,
                        owner_id            UUID             NOT NULL REFERENCES users (id) ON DELETE CASCADE,
                        shared_with_user_id UUID             NOT NULL REFERENCES users (id) ON DELETE CASCADE,
                        permission          share_permission NOT NULL DEFAULT 'viewer',
                        created_at          TIMESTAMPTZ      NOT NULL DEFAULT now(),
                        UNIQUE (resource_type, resource_id, shared_with_user_id)
);

CREATE INDEX shares_shared_with_idx ON shares (shared_with_user_id);
CREATE INDEX shares_resource_idx ON shares (resource_type, resource_id);