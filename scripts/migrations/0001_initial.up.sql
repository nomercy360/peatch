CREATE TABLE users
(
    id              SERIAL PRIMARY KEY,
    first_name      VARCHAR(255),
    last_name       VARCHAR(255),
    chat_id         BIGINT UNIQUE NOT NULL,
    username        VARCHAR(255),
    created_at      TIMESTAMP     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP     NOT NULL DEFAULT NOW(),
    is_published    BOOLEAN       NOT NULL DEFAULT FALSE,
    published_at    TIMESTAMP,
    avatar_url      VARCHAR(512),
    title           VARCHAR(255),
    description     TEXT,
    language        VARCHAR(2)    NOT NULL,
    country         VARCHAR(255),
    city            VARCHAR(255),
    country_code    VARCHAR(2),
    followers_count INTEGER       NOT NULL DEFAULT 0,
    following_count INTEGER       NOT NULL DEFAULT 0,
    requests_count  INTEGER       NOT NULL DEFAULT 0,
    notifications   boolean       NOT NULL DEFAULT true
);

CREATE TABLE badges
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    icon       VARCHAR(255),
    color      VARCHAR(7),
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    user_id    INTEGER REFERENCES users (id)
);

CREATE TABLE opportunities
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    icon        VARCHAR(255),
    color       VARCHAR(7),
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE user_badges
(
    user_id  INTEGER REFERENCES users (id),
    badge_id INTEGER REFERENCES badges (id),
    UNIQUE (user_id, badge_id)
);

-- user_opportunities

CREATE TABLE user_opportunities
(
    user_id        INTEGER REFERENCES users (id),
    opportunity_id INTEGER REFERENCES opportunities (id),
    UNIQUE (user_id, opportunity_id)
);


CREATE TABLE collaborations
(
    id             SERIAL PRIMARY KEY,
    user_id        INTEGER REFERENCES users (id),
    opportunity_id INTEGER REFERENCES opportunities (id),
    title          VARCHAR(255) NOT NULL,
    description    TEXT,
    is_payable     BOOLEAN      NOT NULL DEFAULT FALSE,
    is_published   BOOLEAN      NOT NULL DEFAULT FALSE,
    published_at   TIMESTAMP,
    created_at     TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP    NOT NULL DEFAULT NOW(),
    country        VARCHAR(255),
    city           VARCHAR(255),
    country_code   VARCHAR(2),
    requests_count INTEGER      NOT NULL DEFAULT 0
);


CREATE TABLE user_followers
(
    user_id     INTEGER REFERENCES users (id),
    follower_id INTEGER REFERENCES users (id),
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, follower_id)
);

CREATE OR REPLACE FUNCTION update_follower_counts() RETURNS TRIGGER AS
$$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE users SET followers_count = followers_count + 1 WHERE id = NEW.user_id;
        UPDATE users SET following_count = following_count + 1 WHERE id = NEW.follower_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE users SET followers_count = followers_count - 1 WHERE id = OLD.user_id;
        UPDATE users SET following_count = following_count - 1 WHERE id = OLD.follower_id;
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER update_follower_counts_trigger
    AFTER INSERT OR DELETE
    ON user_followers
    FOR EACH ROW
EXECUTE FUNCTION update_follower_counts();


CREATE TYPE collaboration_request_status AS ENUM ('pending', 'approved', 'rejected');


CREATE TABLE collaboration_requests
(
    id               SERIAL PRIMARY KEY,
    collaboration_id INTEGER REFERENCES collaborations (id),
    user_id          INTEGER REFERENCES users (id),
    message          TEXT,
    created_at       TIMESTAMP                    NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP                    NOT NULL DEFAULT NOW(),
    status           collaboration_request_status NOT NULL DEFAULT 'pending',
    UNIQUE (collaboration_id, user_id)
);

CREATE OR REPLACE FUNCTION update_collaboration_requests_count() RETURNS TRIGGER AS
$$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE collaborations SET requests_count = requests_count + 1 WHERE id = NEW.collaboration_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE collaborations SET requests_count = requests_count - 1 WHERE id = OLD.collaboration_id;
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE user_collaboration_requests
(
    id           SERIAL PRIMARY KEY,
    requester_id INTEGER REFERENCES users (id),
    receiver_id  INTEGER REFERENCES users (id),
    message      TEXT,
    status       collaboration_request_status NOT NULL DEFAULT 'pending',
    created_at   TIMESTAMP                    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP                    NOT NULL DEFAULT NOW(),
    UNIQUE (requester_id, receiver_id)
);

CREATE OR REPLACE FUNCTION update_user_requests_count() RETURNS TRIGGER AS
$$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE users SET requests_count = requests_count + 1 WHERE id = NEW.receiver_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE users SET requests_count = requests_count - 1 WHERE id = OLD.receiver_id;
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;


CREATE table collaboration_badges
(
    collaboration_id INTEGER REFERENCES collaborations (id),
    badge_id         INTEGER REFERENCES badges (id),
    UNIQUE (collaboration_id, badge_id)
);

CREATE table notifications
(
    id         SERIAL PRIMARY KEY,
    user_id    BIGINT,
    message_id VARCHAR(255),
    chat_id    BIGINT,
    text       TEXT,
    image_url  TEXT,
    sent_at    TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
