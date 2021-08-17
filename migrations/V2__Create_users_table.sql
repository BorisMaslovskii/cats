create table users (
    id uuid primary key not null,
    login varchar unique not null,
    password varchar not null,
    admin boolean
);