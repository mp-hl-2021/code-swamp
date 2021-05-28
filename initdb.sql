drop table if exists accounts cascade ;
create table accounts
(
    id        serial primary key,
    login     varchar(255) not null,
    password  varchar(255) not null,
    createdAt timestamp without time zone default now(),
    updatedAt timestamp without time zone default now(),

    unique (login)
);

drop table if exists snippets cascade;
create table snippets
(
    id        serial primary key,
    code      varchar not null,
    uid       int,
    language  varchar(64),
    lifetime  interval not null,
    createdAt timestamp without time zone default now(),
    isChecked bool not null,
    isCorrect bool not null,
    message   varchar not null
);