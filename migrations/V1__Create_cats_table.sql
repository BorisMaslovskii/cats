create table CATS (
    ID serial primary key not null,
    NAME varchar unique not null,
    COLOR varchar not null,
);