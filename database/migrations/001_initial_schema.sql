-- 001_initial_schema.sql

-- Create the service table
CREATE TABLE IF NOT EXISTS service (
    service_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    version TEXT,
    created_at UUID,
    updated_at UUID
);

-- Create the version table
CREATE TABLE IF NOT EXISTS version (
    version_id SERIAL PRIMARY KEY,
    service_id INT NOT NULL,
    version_number INT NOT NULL,
    description TEXT NOT NULL,
    created_at UUID,
    updated_at UUID,
    FOREIGN KEY (service_id) REFERENCES service(service_id)
);
