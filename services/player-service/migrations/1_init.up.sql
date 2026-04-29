CREATE TABLE IF NOT EXISTS players (
    id           VARCHAR(36) PRIMARY KEY,
    name         VARCHAR(50) NOT NULL,
    nickname     VARCHAR(30) NOT NULL UNIQUE,
    position     VARCHAR(20),
    age          INT,
    games_played INT NOT NULL DEFAULT 0,
    wins         INT NOT NULL DEFAULT 0,
    losses       INT NOT NULL DEFAULT 0,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP NOT NULL DEFAULT NOW()
);