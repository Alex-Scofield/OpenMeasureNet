CREATE EXTENSION IF NOT EXISTS postgis; 

CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE versions (
    id SERIAL PRIMARY KEY
);

CREATE TABLE boitiers (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    version_id SERIAL REFERENCES versions(id)
);

CREATE TABLE quantities (
    id SERIAL PRIMARY KEY,
    unit VARCHAR(32) NOT NULL
);

CREATE TABLE observations (
    time TIMESTAMP NOT NULL,
    boitier_id UUID REFERENCES boitiers(id),
    quantity SERIAL REFERENCES quantities(id),
    value REAL NOT NULL,
    location GEOMETRY(POINT, 4326),
    PRIMARY KEY (time, boitier_id, value)
);
