CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS events_db
(
    id          uuid PRIMARY KEY,
    title       TEXT NOT NULL,
    description TEXT
);