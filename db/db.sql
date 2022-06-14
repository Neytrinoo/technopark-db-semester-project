CREATE EXTENSION citext;

-- Tables
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
    id bigserial,
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
    parent_path BIGINT[]    DEFAULT ARRAY []::integer[]
);

CREATE TABLE if not exists Votes
(
    nickname citext  NOT NULL REFERENCES Users (nickname),
    thread   serial  NOT NULL REFERENCES Threads (id),
    voice    integer NOT NULL
);


-- Procedures
CREATE OR REPLACE FUNCTION set_post_parent_path() RETURNS TRIGGER AS
$set_post_parent_path$
BEGIN
    new.parent_path = (SELECT parent_path FROM Posts WHERE id = new.parent) || new.id;
    return new;
END;
$set_post_parent_path$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION add_forum_thread_count() RETURNS TRIGGER AS
$add_forum_thread_count$
BEGIN
    UPDATE Forums SET threads = Forums.threads + 1 WHERE slug = new.forum;
    return new;
END;
$add_forum_thread_count$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION add_forum_posts_count() RETURNS TRIGGER AS
$add_forum_posts_count$
BEGIN
    UPDATE Forums SET posts = Forum.posts + 1 WHERE slug = new.forum;
    return new;
END;
$add_forum_posts_count$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION add_thread_vote() RETURNS TRIGGER AS
$add_thread_vote$
BEGIN
    UPDATE Threads SET votes = Threads.votes + new.voice WHERE id = new.thread;
    return new;
END;
$add_thread_vote$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_thread_vote() RETURNS TRIGGER AS
$update_thread_vote$
BEGIN
    UPDATE Threads SET votes = Threads.votes + (new.voice - old.voice) WHERE id = new.thread;
    return new;
END;
$update_thread_vote$ LANGUAGE plpgsql;

-- Triggers
CREATE TRIGGER set_post_parent_path_trigger
    AFTER INSERT
    ON Posts
    FOR EACH ROW
EXECUTE PROCEDURE set_post_parent_path();
CREATE TRIGGER add_forum_thread_count_trigger
    AFTER INSERT
    ON Threads
    FOR EACH ROW
EXECUTE PROCEDURE add_forum_thread_count();
CREATE TRIGGER add_forum_posts_count_trigger
    AFTER INSERT
    ON Posts
    FOR EACH ROW
EXECUTE PROCEDURE add_forum_posts_count();
CREATE TRIGGER add_thread_vote_trigger
    AFTER INSERT
    ON Votes
    FOR EACH ROW
EXECUTE PROCEDURE add_thread_vote();
CREATE TRIGGER update_thread_vote_trigger
    AFTER UPDATE
    ON Votes
    FOR EACH ROW
EXECUTE PROCEDURE update_thread_vote();


-- Indexes

-- Threads
CREATE INDEX IF NOT EXISTS for_search_users_on_forum_threads ON Threads (forum, author);
CREATE INDEX IF NOT EXISTS for_search_threads_on_forum ON Threads (forum, created);
CREATE INDEX IF NOT EXISTS for_search_by_slug ON Threads USING hash (slug);
CREATE INDEX IF NOT EXISTS for_search_by_forum ON Threads USING hash (forum);

-- Posts
CREATE INDEX IF NOT EXISTS for_search_users_on_forum_posts ON Posts (forum, author);
CREATE INDEX IF NOT EXISTS for_flat_search ON Posts (thread, id);
CREATE INDEX IF NOT EXISTS for_tree_search ON Posts (thread, parent_path, id);
CREATE INDEX IF NOT EXISTS for_parent_tree_search ON Posts ((parent_path[1]), parent_path);
CREATE INDEX IF NOT EXISTS for_search_parents_posts ON Posts (thread, parent, id);

-- User


-- Vote
CREATE UNIQUE INDEX IF NOT EXISTS search_user_vote ON Votes (nickname, thread);