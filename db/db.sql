CREATE EXTENSION citext;

CREATE TABLE if not exists Users
(
    id       bigserial NOT NULL,
    nickname citext    NOT NULL PRIMARY KEY,
    fullname text      NOT NULL,
    about    text,
    email    citext UNIQUE
);

CREATE TABLE if not exists Forums
(
    slug    citext NOT NULL PRIMARY KEY,
    title   text   NOT NULL,
    "user"  citext NOT NULL REFERENCES Users (nickname),
    posts   bigint  DEFAULT 0,
    threads integer DEFAULT 0
);

CREATE TABLE if not exists Threads
(
    id      serial NOT NULL PRIMARY KEY,
    title   text   NOT NULL,
    author  citext NOT NULL REFERENCES Users (nickname),
    forum   citext NOT NULL REFERENCES Forums (slug),
    message text   NOT NULL,
    votes   integer     DEFAULT 0,
    slug    citext,
    created timestamptz DEFAULT now()
);

CREATE TABLE if not exists Posts
(
    id          bigserial NOT NULL PRIMARY KEY,
    parent      bigint      DEFAULT 0,
    author      citext    NOT NULL REFERENCES Users (nickname),
    message     text      NOT NULL,
    isEdited    boolean     DEFAULT false,
    forum       citext    NOT NULL REFERENCES Forums (slug),
    thread      integer REFERENCES Threads (id),
    created     timestamptz DEFAULT now(),
    parent_path BIGINT[]    DEFAULT ARRAY []
);

CREATE TABLE if not exists Votes
(
    nickname citext    NOT NULL REFERENCES Users (nickname),
    thread   serial    NOT NULL REFERENCES Threads (id),
    voice    integer   NOT NULL
);


CREATE OR REPLACE FUNCTION set_post_parent_path() RETURNS TRIGGER AS
$set_post_parent_path$
BEGIN
    new.parent_path = (SELECT parent_path FROM Posts WHERE id = new.parent) || new.id;
    RETURN new;
END;
$set_post_parent_path$ LANGUAGE plpgsql;