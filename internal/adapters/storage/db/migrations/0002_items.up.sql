BEGIN TRANSACTION;

CREATE TYPE item_type AS ENUM ('password', 'card', 'text', 'binary');

CREATE TABLE IF NOT EXISTS items (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    data         BYTEA NOT NULL,
    create_time  TIMESTAMP WITH TIME ZONE NOT NULL,
    update_time  TIMESTAMP WITH TIME ZONE NOT NULL
);

COMMIT;
