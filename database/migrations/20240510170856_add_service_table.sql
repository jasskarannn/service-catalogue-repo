-- +goose Up
-- Create the service table
CREATE TABLE IF NOT EXISTS service (
    service_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    version_count INT,
    created_at UUID
);

-- Create the version table
CREATE TABLE IF NOT EXISTS version (
    version_id SERIAL PRIMARY KEY,
    service_id INT NOT NULL,
    version_number FLOAT NOT NULL,
    description TEXT,
    created_at UUID,
    FOREIGN KEY (service_id) REFERENCES service(service_id)
);

-- -- +goose Down
DROP TABLE IF EXISTS service;

DROP TABLE IF EXISTS version;