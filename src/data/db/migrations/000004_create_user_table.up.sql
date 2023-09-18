-- Active: 1694008398556@@127.0.0.1@5432@greenlight2
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW(),
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    activated bool NOT NULL,
    version INTEGER NOT NULL DEFAULT 1
);