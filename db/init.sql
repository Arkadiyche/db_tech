CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE users
(
    nickname CITEXT PRIMARY KEY NOT NULL,
    email    CITEXT UNIQUE      NOT NULL,
    about    TEXT               NOT NULL,
    fullname TEXT               NOT NULL
);

CREATE UNIQUE INDEX ON users (nickname, email);
CREATE UNIQUE INDEX ON users (nickname, email, about, fullname);
-- CREATE UNIQUE INDEX ON users (nickname ASC); ALREADY EXISTS
CREATE UNIQUE INDEX ON users (nickname DESC);

CREATE UNLOGGED TABLE forums
(
    slug     CITEXT PRIMARY KEY                                   NOT NULL,
    title    TEXT                                                 NOT NULL,
    nickname CITEXT REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    posts    INTEGER DEFAULT 0                                    NOT NULL,
    threads  INTEGER DEFAULT 0                                    NOT NULL
);



CREATE UNLOGGED TABLE forum_users
(
    author CITEXT REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    slug   CITEXT REFERENCES forums (slug) ON DELETE CASCADE    NOT NULL,
    PRIMARY KEY (slug, author)
);

CREATE INDEX ON forum_users (slug);
CREATE INDEX ON forum_users (author);

CREATE UNLOGGED TABLE threads
(
    id         SERIAL PRIMARY KEY                                    NOT NULL,
    author     CITEXT REFERENCES users (nickname) ON DELETE CASCADE  NOT NULL,
    created    TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    forum_slug CITEXT REFERENCES forums (slug) ON DELETE CASCADE     NOT NULL,
    message    TEXT                                                  NOT NULL,
    slug       CITEXT,
    title      TEXT                                                  NOT NULL,
    votes      INTEGER                     DEFAULT 0                 NOT NULL
);

CREATE INDEX ON threads(slug, author);
CREATE INDEX ON threads(forum_slug, created ASC);
CREATE INDEX ON threads(forum_slug, created DESC);
-- CREATE INDEX ON threads(id); ALREADY EXISTS
CREATE INDEX ON threads(slug, id);
CREATE INDEX ON threads(id, forum_slug);
CREATE INDEX ON threads(slug, id, forum_slug);


CREATE UNLOGGED TABLE posts
(
    author     CITEXT REFERENCES users (nickname) ON DELETE CASCADE  NOT NULL,
    created    TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    forum_slug CITEXT REFERENCES forums (slug) ON DELETE CASCADE     NOT NULL,
    id         SERIAL PRIMARY KEY                                    NOT NULL,
    edited     BOOL                        DEFAULT 'false'           NOT NULL,
    message    TEXT                                                  NOT NULL,
    parent     INTEGER                                               NOT NULL,
    thread     INTEGER REFERENCES threads (id) ON DELETE CASCADE     NOT NULL,
    path       INTEGER ARRAY               DEFAULT '{}'              NOT NULL
);

CREATE UNIQUE INDEX ON posts(id, thread);
CREATE UNIQUE INDEX ON posts(id, author);
CREATE INDEX ON posts(thread, path DESC);
CREATE INDEX ON posts(thread, path ASC);
CREATE INDEX ON posts(thread, id DESC);
CREATE INDEX ON posts(thread, id ASC);

CREATE UNLOGGED TABLE votes
(
    nickname  CITEXT REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    thread_id INTEGER REFERENCES threads (id) ON DELETE CASCADE    NOT NULL,
    vote      SMALLINT                                             NOT NULL,
    PRIMARY KEY (thread_id, nickname)
);

CREATE UNIQUE INDEX ON votes(thread_id, nickname);

CREATE INDEX ON threads (slug, id);

CREATE FUNCTION update_path() RETURNS TRIGGER AS
$$
DECLARE
    temp INT ARRAY;
    t    INTEGER;
BEGIN
    IF new.parent ISNULL OR new.parent = 0 THEN
        new.path = ARRAY [new.id];
    ELSE
        SELECT thread
        INTO t
        FROM posts
        WHERE id = new.parent;
        IF t ISNULL OR t <> new.thread THEN
            RAISE EXCEPTION 'Not in this thread ID ' USING HINT = 'Please check your parent ID';
        END IF;

        SELECT path
        INTO temp
        FROM posts
        WHERE id = new.parent;
        new.path = array_append(temp, new.id);

    END IF;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_posts_path
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_path();

-- VOTE VALUE UPDATE
CREATE FUNCTION vote_count_upd() RETURNS TRIGGER AS
$$
BEGIN
    IF (old.vote != new.vote) THEN
        UPDATE threads
        SET votes = (votes - old.vote + new.vote)
        WHERE id = new.thread_id;
    END IF;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER vote_count_upd
    AFTER UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE vote_count_upd();

CREATE FUNCTION  vote_count_insert() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE threads
    SET votes = (votes + new.vote)
    WHERE id = new.thread_id;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER vote_count_insert
    AFTER INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE vote_count_insert();


-- UPDATE FORUM_USER TABLE AFTER INSERTS
CREATE FUNCTION  insert_forum_user_from_threads_or_psoots() RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO forum_users
    VALUES (new.author, new.forum_slug)
    ON CONFLICT DO NOTHING;
    RETURN NULL;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER update_forum_user_from_threads
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE insert_forum_user_from_threads_or_psoots();

CREATE TRIGGER update_forum_user_from_posts
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE insert_forum_user_from_threads_or_psoots();


CREATE FUNCTION  update_forum_counter_posts() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE forums
    SET posts = posts + 1
    WHERE slug = new.forum_slug;

    RETURN NULL;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER update_forum_counters_after_post_insert
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_forum_counter_posts();

CREATE FUNCTION  update_forum_counter_threads() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE forums
    SET threads = threads + 1
    WHERE slug = new.forum_slug;

    RETURN NULL;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER update_forum_counters_after_thread_insert
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE update_forum_counter_threads();