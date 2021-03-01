CREATE TABLE hours (
    id INTEGER PRIMARY KEY,
    date text not null,
    ticket varchar(50) not null,
    title varchar(150) not null,
    comment text,
    hours float not null
);

CREATE TABLE hours_templates (
    id INTEGER PRIMARY KEY,
    ticket varchar(50) not null,
    title varchar(150) not null,
    comment text,
    hours float
);