CREATE TABLE IF NOT EXISTS dialogs (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    state TEXT NOT NULL DEFAULT 'STARTED', -- STARTED, FINISHED, CANCELLED
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    command TEXT NOT NULL, -- create_notification, add_medicine, etc.
    context TEXT, -- store intermediate data
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
