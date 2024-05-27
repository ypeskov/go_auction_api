CREATE TABLE item_comments (
    id SERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL REFERENCES items(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    comment TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);