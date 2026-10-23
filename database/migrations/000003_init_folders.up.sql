CREATE TABLE folders (
                         id         UUID PRIMARY KEY,
                         owner_id   UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
                         parent_id  UUID        REFERENCES folders (id) ON DELETE CASCADE,
                         name       TEXT        NOT NULL,
                         created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                         updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                         deleted_at TIMESTAMPTZ,
                         CONSTRAINT folders_name_not_empty CHECK (length(btrim(name)) > 0),
                         CONSTRAINT folders_no_self_parent CHECK (parent_id IS DISTINCT FROM id)
    );


CREATE UNIQUE INDEX folders_unique_name_in_parent
    ON folders (owner_id, parent_id, name)
    WHERE deleted_at IS NULL AND parent_id IS NOT NULL;

CREATE UNIQUE INDEX folders_unique_name_in_root
    ON folders (owner_id, name)
    WHERE deleted_at IS NULL AND parent_id IS NULL;

CREATE INDEX folders_owner_parent_idx ON folders (owner_id, parent_id);
CREATE INDEX folders_deleted_at_idx ON folders (deleted_at) WHERE deleted_at IS NOT NULL;