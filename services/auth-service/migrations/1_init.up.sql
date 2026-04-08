CREATE TABLE IF NOT EXISTS users (
    id           VARCHAR(36) PRIMARY KEY,
    name         VARCHAR(50) NOT NULL,
    nickname     VARCHAR(30) NOT NULL UNIQUE,
    email        VARCHAR(255) UNIQUE,
    phone        VARCHAR(20) UNIQUE,
    google_id    VARCHAR(255) UNIQUE,
    pass_hash    BYTEA,
    role         VARCHAR(20) NOT NULL DEFAULT 'player',
    trust_rating INT NOT NULL DEFAULT 100,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP NOT NULL DEFAULT NOW()
)