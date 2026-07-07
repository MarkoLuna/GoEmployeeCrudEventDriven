CREATE TABLE IF NOT EXISTS clients (
    id              TEXT PRIMARY KEY,
    client_id       TEXT UNIQUE NOT NULL,
    client_secret   TEXT NOT NULL,
    enabled         BOOLEAN DEFAULT true,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
