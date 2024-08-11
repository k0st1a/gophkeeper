BEGIN TRANSACTION;

CREATE TYPE item_type AS ENUM ('login', 'card', 'text', 'binary');

CREATE TABLE IF NOT EXISTS items (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    name        TEXT NOT NULL,
    type        item_type NOT NULL,
    data        BYTEA NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    UNIQUE      (name, user_id)
);

COMMIT;
