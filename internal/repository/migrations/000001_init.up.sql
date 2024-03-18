CREATE TABLE IF NOT EXISTS actors
(
    id            serial not null unique,
    name          varchar(255),
    gender        varchar(10),
    birthday DATE
);

CREATE TABLE IF NOT EXISTS films
(
    id           serial not null unique,
    title        varchar(255),
    description  varchar(1000),
    release_date date,
    rating       decimal(4, 2) check ( rating >= 0 and rating <= 10 )
);

CREATE TABLE IF NOT EXISTS actors_films
(
    id       serial                                       not null unique,
    actor_id int references actors (id) on delete cascade not null,
    film_id  int references films (id) on delete cascade  not null
);

CREATE TABLE IF NOT EXISTS users
(
    id       serial       not null unique,
    username varchar(255) not null ,
    role     varchar(5) not null
);