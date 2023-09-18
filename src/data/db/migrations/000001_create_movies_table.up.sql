-- Active: 1694008398556@@127.0.0.1@5432@greenlight2
CREATE TABLE IF NOT EXISTS movies (
    id bigserial PRIMARY KEY,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    title TEXT NOT NULL,
    year INTEGER NOT NULL,
    runtime INTEGER NOT NULL,
    genres TEXT[] NOT NULL,
    version INTEGER NULL DEFAULT 1
);


