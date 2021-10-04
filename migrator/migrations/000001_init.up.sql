CREATE TABLE users
(
    id            uuid PRIMARY KEY NOT NULL UNIQUE,
    email         varchar(255)     NOT NULL UNIQUE,
    password_hash varchar(255)     NOT NULL UNIQUE
);

CREATE TABLE links
(
    id        uuid PRIMARY KEY                             NOT NULL UNIQUE,
    user_id   uuid REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    link      varchar(4096)                                NOT NULL UNIQUE,
    link_hash varchar(40)                                  NOT NULL UNIQUE
);