CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    nickname citext COLLATE "ucs_basic" PRIMARY KEY,
    fullname VARCHAR NOT NULL,
    email citext NOT NULL UNIQUE,
    about VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS forums (
    slug citext PRIMARY KEY,
    title VARCHAR NOT NULL,
    user_nick citext REFERENCES users(nickname),
    posts BIGINT DEFAULT 0,
    threads BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS forum_users (
    forum citext REFERENCES forums(slug) ON DELETE CASCADE,
    user_nick citext REFERENCES users(nickname) ON DELETE CASCADE,
    PRIMARY KEY (user_nick, forum)
);

CREATE TABLE IF NOT EXISTS threads (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    title VARCHAR NOT NULL,
    author citext NOT NULL REFERENCES users(nickname),
    forum citext NOT NULL REFERENCES forums(slug) ON DELETE CASCADE,
    message TEXT NOT NULL,
    slug citext NOT NULL,
    votes INT NOT NULL DEFAULT 0,
    post_tree BIGINT[] DEFAULT '{}',
    created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    parent BIGINT DEFAULT 0,
    author citext REFERENCES users(nickname),
    forum citext REFERENCES forums(slug) ON DELETE CASCADE,
    thread BIGSERIAL REFERENCES threads(id) ON DELETE CASCADE,
    is_edited BOOLEAN DEFAULT FALSE,
    message TEXT NOT NULL,
    created TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    path BIGINT[] DEFAULT ARRAY []::BIGINT[]
);

CREATE TABLE IF NOT EXISTS votes
(
    nickname  citext REFERENCES users (nickname) ON DELETE NO ACTION NOT NULL,
    thread BIGSERIAL REFERENCES threads (id) ON DELETE CASCADE    NOT NULL,
    voice     SMALLINT CHECK ( voice BETWEEN -1 AND 1 )              NOT NULL,
    PRIMARY KEY (nickname, thread)
);


-- Functions and triggers for updating votes in threads
CREATE OR REPLACE FUNCTION insert_trigger_thread_votes() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE threads SET votes = votes + new.voice WHERE id = new.thread;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_trigger_thread_votes() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE threads SET votes = votes + new.voice - old.voice WHERE id = new.thread;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_trigger_thread_votes
    AFTER INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE insert_trigger_thread_votes();

CREATE TRIGGER update_trigger_thread_votes
    AFTER UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE update_trigger_thread_votes();

CREATE OR REPLACE FUNCTION update_path_trigger() RETURNS TRIGGER AS
$$
BEGIN
    new.path = (SELECT path FROM posts WHERE id = new.parent) || new.id;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_path_trigger
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_path_trigger();

-- Functions and triggers for counting posts and threads in forums
CREATE OR REPLACE FUNCTION insert_trigger_forum_posts() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE forums SET posts = posts + 1 WHERE slug = new.forum;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_trigger_forum_posts
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE insert_trigger_forum_posts();

CREATE OR REPLACE FUNCTION insert_trigger_forum_threads() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE forums SET threads = threads + 1 WHERE slug = new.forum;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_trigger_forum_threads
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE insert_trigger_forum_threads();
