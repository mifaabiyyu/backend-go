CREATE TABLE IF NOT EXISTS users(
    "id" bigserial PRIMARY KEY,
    "email" varchar UNIQUE NOT NULL,
    "username" varchar NOT NULL,
    "full_name" varchar NOT NULL,
    "password" varchar NOT NULL,
    "verified" boolean NOT NULL DEFAULT FALSE,
    "verified_at" timestamptz,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);