CREATE TABLE modules (
    id         TEXT      NOT NULL PRIMARY KEY,
    file_path  TEXT      NOT NULL,
    enabled    BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_module_id ON modules( id );
CREATE INDEX IF NOT EXISTS idx_enabled ON modules( enabled );
