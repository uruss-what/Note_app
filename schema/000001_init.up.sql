--Active: 1703365526307@@127.0.0.1@5433@postgres
CREATE SCHEMA IF NOT EXISTS db;

CREATE TABLE users
(
    id SERIAL NOT NULL UNIQUE,
    name VARCHAR NOT NULL,
    username VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL
);

CREATE TABLE todo_lists
(
    id SERIAL NOT NULL UNIQUE,
    title VARCHAR NOT NULL,
    description VARCHAR
);

CREATE TABLE users_lists
(
    id SERIAL NOT NULL UNIQUE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    list_id INT REFERENCES todo_lists(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE todo_items
(
    id SERIAL NOT NULL UNIQUE,
    title VARCHAR NOT NULL,
    description VARCHAR,
    done BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE lists_items
(
    id SERIAL NOT NULL UNIQUE,
    item_id INT REFERENCES todo_items(id) ON DELETE CASCADE NOT NULL,
    list_id INT REFERENCES todo_lists(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE roles
(
    id SERIAL NOT NULL UNIQUE,
    title VARCHAR NOT NULL,
    description VARCHAR
);

CREATE TABLE users_roles
(
    id SERIAL NOT NULL UNIQUE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    role_id INT REFERENCES roles(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE statuses
(
    id SERIAL NOT NULL UNIQUE,
    description VARCHAR
);

CREATE TABLE users_statuses
(
    id SERIAL NOT NULL UNIQUE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    status_id INT REFERENCES statuses(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE comments
(
    id SERIAL NOT NULL UNIQUE,
    description VARCHAR,
    author_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    receiver_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL
);

--CREATE USER urusswhat WITH PASSWORD 'fpeople';

--ALTER USER urusswhat WITH SUPERUSER;