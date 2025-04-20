CREATE TABLE users (
    id bigserial PRIMARY KEY,
    email TEXT NOT NULL,
    username TEXT NOT NULL,
    full_name TEXT NOT NULL,
    password TEXT NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    verified_at timestamptz NOT NULL DEFAULT (now()),
    created_at timestamptz NOT NULL DEFAULT (now()),
    updated_at timestamptz NOT NULL DEFAULT (now()),
    role_id INT REFERENCES roles(id)
);