BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS users (
    id         BIGSERIAL PRIMARY KEY,
    login      TEXT UNIQUE NOT NULL,
    password   TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

COMMIT;
