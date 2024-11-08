-- 000004_create_events_and_chats.up.sql
BEGIN;

CREATE TABLE IF NOT EXISTS events (
    event_id VARCHAR(10) NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    cover_image_url VARCHAR(255),
    location VARCHAR(100),
    organized_by VARCHAR(100),
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    created_by VARCHAR(10) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (created_by) REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS chats (
    chat_id VARCHAR(10) NOT NULL PRIMARY KEY,
    sender VARCHAR(10) NOT NULL,
    receiver VARCHAR(10) NOT NULL,
    type VARCHAR(10) NOT NULL, -- e.g., text, image
    text_message VARCHAR(255),
    image_url_message VARCHAR(255),
    timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender) REFERENCES users(user_id),
    FOREIGN KEY (receiver) REFERENCES users(user_id)
);

COMMIT;