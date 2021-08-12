create table cats (
    id serial primary key not null,
    name varchar unique not null,
    color varchar not null
);