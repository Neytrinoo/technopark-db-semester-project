CREATE EXTENSION IF NOT EXISTS citext;

-- Tables
CREATE UNLOGGED TABLE if not exists Users
(
    id       bigserial          NOT NULL,
    nickname citext COLLATE "C" NOT NULL PRIMARY KEY, -- для побайтового сравнения в нижнем регистре добавляем COLLATE "C"
    fullname text               NOT NULL,
    about    text,
    email    citext             NOT NULL UNIQUE
);

CREATE UNLOGGED TABLE if not exists Forums
(
    id      bigserial,
    slug    citext             NOT NULL PRIMARY KEY,
    title   text               NOT NULL,
    "user"  citext COLLATE "C" NOT NULL REFERENCES Users (nickname),
    posts   bigint DEFAULT 0,
    threads bigint DEFAULT 0
);

CREATE UNLOGGED TABLE if not exists Threads
(
    id      bigserial          NOT NULL PRIMARY KEY,
    title   text               NOT NULL,
    author  citext COLLATE "C" NOT NULL REFERENCES Users (nickname),
    forum   citext             NOT NULL REFERENCES Forums (slug),
    message text               NOT NULL,
    votes   integer                  DEFAULT 0,
    slug    citext             NOT NULL,
    created timestamptz DEFAULT now()
);

CREATE UNLOGGED TABLE if not exists Posts
(
    id          bigserial          NOT NULL PRIMARY KEY,
    parent      integer                  DEFAULT 0,
    author      citext COLLATE "C" NOT NULL REFERENCES Users (nickname),
    message     text               NOT NULL,
    isEdited    boolean                  DEFAULT false,
    forum       citext             NOT NULL REFERENCES Forums (slug),
    thread      integer REFERENCES Threads (id),
    created     timestamptz DEFAULT now(),
    parent_path BIGINT[]                 DEFAULT ARRAY []::integer[]
);

CREATE UNLOGGED TABLE IF NOT EXISTS ForumUsers
(
    nickname citext COLLATE "C" NOT NULL REFERENCES Users (nickname),
    fullname text               NOT NULL,
    about    text,
    email    citext             NOT NULL,
    forum    citext             NOT NULL REFERENCES Forums (slug),
    CONSTRAINT nickname_forum_unique UNIQUE (nickname, forum)
);

CREATE UNLOGGED TABLE if not exists Votes
(
    nickname citext COLLATE "C" NOT NULL REFERENCES Users (nickname),
    thread   serial             NOT NULL REFERENCES Threads (id),
    voice    integer            NOT NULL
);


-- Procedures

CREATE OR REPLACE FUNCTION forum_users_update() RETURNS TRIGGER AS
$forum_users_update$
DECLARE
    nickname_param citext;
    fullname_param text;
    about_param    text;
    email_param    citext;
BEGIN
    SELECT t.nickname, t.fullname, t.about, t.email
    FROM Users AS t
    WHERE t.nickname = new.author
    INTO nickname_param, fullname_param, about_param, email_param;

    INSERT INTO ForumUsers (nickname, fullname, about, email, forum)
    VALUES (nickname_param, fullname_param, about_param, email_param, new.forum)
    ON CONFLICT DO NOTHING;

    return new;
END;
$forum_users_update$ LANGUAGE plpgsql;


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
    UPDATE Forums SET posts = Forums.posts + 1 WHERE slug = new.forum;
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
CREATE TRIGGER forum_users_for_post
    AFTER INSERT
    ON Posts
    FOR EACH ROW
EXECUTE PROCEDURE forum_users_update();

CREATE TRIGGER forum_users_for_thread
    AFTER INSERT
    ON Threads
    FOR EACH ROW
EXECUTE PROCEDURE forum_users_update();

CREATE TRIGGER set_post_parent_path_trigger
    BEFORE INSERT
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
CREATE INDEX IF NOT EXISTS thread_id_hash ON Threads using hash (id);

-- Posts
CREATE INDEX IF NOT EXISTS for_search_users_on_forum_posts ON Posts (forum, author);
CREATE INDEX IF NOT EXISTS for_flat_search ON Posts (thread, id);
CREATE INDEX IF NOT EXISTS for_tree_search ON Posts (thread, parent_path, id);
CREATE INDEX IF NOT EXISTS for_parent_tree_search ON Posts ((parent_path[1]), parent_path);
CREATE INDEX IF NOT EXISTS for_search_parents_posts ON Posts (thread, parent, id);
CREATE INDEX IF NOT EXISTS post_id_hash ON Posts using hash (id);
CREATE INDEX IF NOT EXISTS post_thread_hash ON Posts using hash (thread);

-- User
CREATE INDEX IF NOT EXISTS user_nickname_hash ON Users using hash (nickname);
CREATE INDEX IF NOT EXISTS user_email_hash ON Users using hash (email);

-- Vote
CREATE UNIQUE INDEX IF NOT EXISTS search_user_vote ON Votes (nickname, thread, voice);
--CREATE UNIQUE INDEX IF NOT EXISTS search_user_vote_hash ON Votes using hash (nickname, thread);

-- Forum
CREATE INDEX IF NOT EXISTS forum_slug_hash ON Forums using hash (slug);

-- ForumUsers
CREATE INDEX IF NOT EXISTS forum_users_forum ON ForumUsers (forum, nickname);

VACUUM ANALYZE;
