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
    CHECK (LENGTH(title) > 1 AND rating > 0.00 AND rating < 10.00)
);

CREATE TABLE IF NOT EXISTS actors_to_films (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id UUID NOT NULL,
    film_id UUID NOT NULL,
    CONSTRAINT "FK_actor_id" FOREIGN KEY ("actor_id") REFERENCES "actors" ("id"),
    CONSTRAINT "FK_film_id" FOREIGN KEY ("film_id") REFERENCES "films" ("id")
);

CREATE UNIQUE INDEX "actors_to_films_actor_id_film_id"
    ON "actors_to_films"
    USING btree
    ("actor_id", "film_id");

