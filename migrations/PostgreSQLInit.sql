CREATE TYPE users_role AS enum ('Admin', 'User');

CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role users_role NOT NULL DEFAULT 'User',
    CHECK (LENGTH(login) > 3 and LENGTH(password) > 3)
);

CREATE TYPE actors_gender AS enum ('Male', 'Female');

CREATE TABLE IF NOT EXISTS actors(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    gender actors_gender NOT NULL,
    birth_date DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS films(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(150),
    description VARCHAR(1000),
    release_date DATE,
    rating DECIMAL(4, 2),
    CHECK (LENGTH(title) > 1 AND rating >= 0.00 AND rating <= 10.00)
);

CREATE TABLE IF NOT EXISTS actors_to_films (
    actor_id UUID NOT NULL REFERENCES actors(id) ON DELETE CASCADE,
    film_id UUID NOT NULL REFERENCES films(id) ON DELETE CASCADE,
    CONSTRAINT actors_to_films_pk PRIMARY KEY(actor_id,film_id)
);

CREATE INDEX idx_gin_films
ON films
USING gin(title);

CREATE INDEX idx_gin_actors
ON actors
USING gin(name);

