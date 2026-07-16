CREATE EXTENSION IF NOT EXISTS postgis; 

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE invite_codes (
    id SERIAL PRIMARY KEY,
    code VARCHAR(32) UNIQUE NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    used_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE versions (
    id SERIAL PRIMARY KEY,
    release_date TIMESTAMP NOT NULL
);

CREATE TABLE nodes (
    id SERIAL PRIMARY KEY,
    user_id SERIAL REFERENCES users(id),
    version_id SERIAL REFERENCES versions(id),
    password VARCHAR(255),
    dashboard_uid VARCHAR(64) UNIQUE
);

CREATE TABLE quantities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    unit VARCHAR(32) NOT NULL
);

CREATE TABLE observations (
    time TIMESTAMP NOT NULL,
    node_id SERIAL REFERENCES nodes(id),
    quantity SERIAL REFERENCES quantities(id),
    value REAL NOT NULL,
    location GEOMETRY(POINT, 4326),
    PRIMARY KEY (time, node_id, quantity)
);
