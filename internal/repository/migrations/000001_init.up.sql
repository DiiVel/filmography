CREATE TABLE actors
(
    id serial not null unique,
    name VARCHAR(255),
    gender VARCHAR(10),
    date_of_birth DATE
);

CREATE TABLE films
(
    id serial not null unique,
    title varchar(255),
    description varchar(1000),
    release_date date,
    rating decimal(4, 2) check ( rating >= 0 and rating <= 10 )
);

CREATE TABLE actors_films
(
    id       serial                                       not null unique,
    actor_id int references actors (id) on delete cascade not null,
    film_id  int references films (id) on delete cascade  not null
);