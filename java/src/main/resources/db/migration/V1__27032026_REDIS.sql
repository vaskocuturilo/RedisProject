CREATE TABLE IF NOT EXISTS events_db
(
    id          VARCHAR(36) PRIMARY KEY,
    title       TEXT NOT NULL,
    description TEXT
);