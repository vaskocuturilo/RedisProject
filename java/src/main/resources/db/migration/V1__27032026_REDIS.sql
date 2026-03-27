CREATE TABLE events
(
    id          VARCHAR(36) PRIMARY KEY,
    title       TEXT NOT NULL,
    description TEXT
);

CREATE TABLE users
(
    id   VARCHAR(36) PRIMARY KEY,
    name TEXT    NOT NULL,
    age  INTEGER NOT NULL
);

CREATE TABLE user_events
(
    user_id  VARCHAR(36) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    event_id VARCHAR(36) NOT NULL REFERENCES events (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, event_id)
);