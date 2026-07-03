CREATE EXTENSION IF NOT EXISTS postgis; 

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE versions (
    id SERIAL PRIMARY KEY,
    release_date TIMESTAMP NOT NULL
);

CREATE TABLE boitiers (
    id SERIAL PRIMARY KEY,
    user_id SERIAL REFERENCES users(id),
    version_id SERIAL REFERENCES versions(id)
);

CREATE TABLE quantities (
    id SERIAL PRIMARY KEY,
    unit VARCHAR(32) NOT NULL
);

CREATE TABLE observations (
    time TIMESTAMP NOT NULL,
    boitier_id SERIAL REFERENCES boitiers(id),
    quantity SERIAL REFERENCES quantities(id),
    value REAL NOT NULL,
    location GEOMETRY(POINT, 4326),
    PRIMARY KEY (time, boitier_id, value)
);
