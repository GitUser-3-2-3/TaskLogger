CREATE TABLE users
(
    user_id       integer primary key,
    username      varchar(50)  not null unique,
    email         varchar(100) not null unique,
    firstname     varchar(50)  not null,
    lastname      varchar(50)  not null,
    password_hash varchar(255) not null,
    created_at    datetime     not null default current_timestamp
);
