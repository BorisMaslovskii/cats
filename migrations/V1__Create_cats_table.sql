create table cats (
    id uuid primary key not null,
    name varchar unique not null,
    color varchar not null
);