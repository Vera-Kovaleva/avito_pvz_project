CREATE TYPE city AS ENUM ('Москва', 'Санкт-Петербург', 'Казань');

CREATE TABLE IF NOT EXISTS pvz (
    id UUID PRIMARY KEY,
    city city NOT NULL,
    registered_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TYPE status AS ENUM ('in_progress', 'close');

CREATE TABLE IF NOT EXISTS receptions (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    pvz_id UUID NOT NULL,
    status status NOT NULL,
    FOREIGN KEY(pvz_id) REFERENCES pvz(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX reception_in_progress_unique ON receptions (id, status) WHERE status = 'in_progress';


CREATE TYPE product_type AS ENUM ('электроника', 'одежда', 'обувь');

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    type product_type NOT NULL,
    reception_id UUID NOT NULL,
    FOREIGN KEY(reception_id) REFERENCES receptions(id) ON DELETE CASCADE
);

CREATE TYPE role AS ENUM ('employee', 'moderator');

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL,
    role role NOT NULL,
    password_hash TEXT NOT NULL,
    token TEXT NOT NULL
);