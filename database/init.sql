CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

INSERT INTO users (name, email) VALUES
    ('Max Mustermann', 'max@example.com'),
    ('Erika Musterfrau', 'erika@example.com'),
    ('John Doe', 'john@example.com')
ON CONFLICT (email) DO NOTHING;