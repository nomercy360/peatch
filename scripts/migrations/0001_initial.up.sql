CREATE TABLE users
(
    id                       SERIAL PRIMARY KEY,
    first_name               VARCHAR(255),
    last_name                VARCHAR(255),
    chat_id                  BIGINT UNIQUE NOT NULL,
    username                 VARCHAR(255),
    created_at               timestamptz   NOT NULL DEFAULT NOW(),
    updated_at               timestamptz   NOT NULL DEFAULT NOW(),
    published_at             timestamptz,
    hidden_at                timestamptz            DEFAULT NOW(),
    notifications_enabled_at timestamptz,
    avatar_url               VARCHAR(512),
    title                    VARCHAR(255),
    description              TEXT,
    likes_count              INTEGER       NOT NULL DEFAULT 0,
    language_code            VARCHAR(2)    NOT NULL DEFAULT 'en',
    country                  VARCHAR(255),
    city                     VARCHAR(255),
    country_code             VARCHAR(2),
    followers_count          INTEGER       NOT NULL DEFAULT 0,
    following_count          INTEGER       NOT NULL DEFAULT 0,
    requests_count           INTEGER       NOT NULL DEFAULT 0,
    review_status            VARCHAR(255)  NOT NULL DEFAULT 'pending',
    peatch_points            INTEGER       NOT NULL DEFAULT 0,
    referrer_id              INTEGER REFERENCES users (id),
    last_check_in            timestamptz
);

CREATE OR REPLACE FUNCTION update_referrer_points()
    RETURNS TRIGGER AS
$$
BEGIN
    IF NEW.referrer_id IS NOT NULL THEN
        UPDATE users
        SET peatch_points = peatch_points + 100
        WHERE id = NEW.referrer_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_referrer_points_trigger
    AFTER INSERT
    ON users
    FOR EACH ROW
EXECUTE FUNCTION update_referrer_points();

CREATE INDEX users_chat_id_index ON users (chat_id);
CREATE UNIQUE INDEX users_username_index ON users (username);
CREATE INDEX review_status_index ON users (review_status);

CREATE TABLE user_interactions
(
    id               SERIAL PRIMARY KEY,
    user_id          INTEGER REFERENCES users (id),
    target_user_id   INTEGER REFERENCES users (id),
    interaction_type VARCHAR(10),
    created_at       timestamptz DEFAULT NOW()
);

CREATE TABLE badges
(
    id         SERIAL PRIMARY KEY,
    text       VARCHAR(255) NOT NULL,
    icon       VARCHAR(255),
    color      VARCHAR(7),
    created_at timestamptz DEFAULT NOW()
);

CREATE TABLE opportunities
(
    id          SERIAL PRIMARY KEY,
    text        VARCHAR(255) NOT NULL,
    description TEXT,
    icon        VARCHAR(255),
    color       VARCHAR(7),
    created_at  timestamptz DEFAULT NOW()
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
    published_at   timestamptz,
    created_at     timestamptz  NOT NULL DEFAULT NOW(),
    updated_at     timestamptz  NOT NULL DEFAULT NOW(),
    hidden_at      timestamptz,
    country        VARCHAR(255),
    city           VARCHAR(255),
    country_code   VARCHAR(2),
    likes_count    INTEGER      NOT NULL DEFAULT 0,
    requests_count INTEGER      NOT NULL DEFAULT 0
);


CREATE TABLE user_followers
(
    user_id     INTEGER REFERENCES users (id),
    follower_id INTEGER REFERENCES users (id),
    created_at  timestamptz DEFAULT NOW(),
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

CREATE OR REPLACE FUNCTION update_collaborations_hidden_at() RETURNS TRIGGER AS
$$
BEGIN
    IF NEW.hidden_at IS NOT NULL THEN
        UPDATE collaborations
        SET hidden_at = NOW()
        WHERE user_id = NEW.id
          AND hidden_at IS NULL;
    ELSE
        UPDATE collaborations
        SET hidden_at = NULL
        WHERE user_id = NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_collaborations_hidden_at
    AFTER UPDATE OF hidden_at
    ON users
    FOR EACH ROW
    WHEN (OLD.hidden_at IS DISTINCT FROM NEW.hidden_at)
EXECUTE FUNCTION update_collaborations_hidden_at();

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
    created_at       timestamptz                  NOT NULL DEFAULT NOW(),
    updated_at       timestamptz                  NOT NULL DEFAULT NOW(),
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
    user_id      INTEGER REFERENCES users (id),
    requester_id INTEGER REFERENCES users (id),
    message      TEXT,
    status       collaboration_request_status NOT NULL DEFAULT 'pending',
    created_at   timestamptz                  NOT NULL DEFAULT NOW(),
    updated_at   timestamptz                  NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, requester_id)
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
    id                SERIAL PRIMARY KEY,
    user_id           BIGINT,
    message_id        VARCHAR(255),
    chat_id           BIGINT,
    sent_at           timestamptz,
    created_at        timestamptz NOT NULL DEFAULT NOW(),
    notification_type VARCHAR(255),
    entity_id         INTEGER,
    entity_type       VARCHAR(255),
    UNIQUE (user_id, entity_id, entity_type)
);

create table locations
(
    id           serial
        primary key,
    country_code char(2)      not null,
    country      varchar(100) not null,
    city         varchar(100) not null,
    population   bigint       not null default 0
);

CREATE TABLE posts
(
    id           SERIAL PRIMARY KEY,
    user_id      INTEGER REFERENCES users (id),
    title        TEXT,
    description  TEXT,
    image_url    TEXT,
    country      VARCHAR(255),
    city         VARCHAR(255),
    country_code VARCHAR(2),
    likes_count  INTEGER     NOT NULL DEFAULT 0,
    hidden_at    timestamptz,
    created_at   timestamptz NOT NULL DEFAULT NOW(),
    updated_at   timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE likes
(
    id           SERIAL PRIMARY KEY,
    user_id      INTEGER REFERENCES users (id),
    content_id   INTEGER     NOT NULL,
    content_type VARCHAR(20) NOT NULL,
    created_at   timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, content_id, content_type)
);

CREATE INDEX likes_user_id_index ON likes (user_id);
CREATE INDEX likes_content_id_index ON likes (content_id);
CREATE INDEX likes_content_type_index ON likes (content_type);


CREATE OR REPLACE FUNCTION update_peatch_points()
    RETURNS TRIGGER AS
$$
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.content_type = 'user' THEN
            UPDATE users
            SET peatch_points = peatch_points + 1,
                likes_count   = likes_count + 1
            WHERE id = NEW.content_id;
        ELSIF NEW.content_type = 'post' THEN
            UPDATE users
            SET peatch_points = peatch_points + 1
            WHERE id = (SELECT user_id FROM posts WHERE id = NEW.content_id);
            UPDATE posts
            SET likes_count = likes_count + 1
            WHERE id = NEW.content_id;
        ELSIF NEW.content_type = 'collaboration' THEN
            UPDATE users
            SET peatch_points = peatch_points + 1
            WHERE id = (SELECT user_id FROM collaborations WHERE id = NEW.content_id);
            UPDATE collaborations
            SET likes_count = likes_count + 1
            WHERE id = NEW.content_id;
        END IF;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        IF OLD.content_type = 'user' THEN
            UPDATE users
            SET peatch_points = peatch_points - 1,
                likes_count   = likes_count - 1
            WHERE id = OLD.content_id;
        ELSIF OLD.content_type = 'post' THEN
            UPDATE users
            SET peatch_points = peatch_points - 1
            WHERE id = (SELECT user_id FROM posts WHERE id = OLD.content_id);
            UPDATE posts
            SET likes_count = likes_count - 1
            WHERE id = OLD.content_id;
        ELSIF OLD.content_type = 'collaboration' THEN
            UPDATE users
            SET peatch_points = peatch_points - 1
            WHERE id = (SELECT user_id FROM collaborations WHERE id = OLD.content_id);
            UPDATE collaborations
            SET likes_count = likes_count - 1
            WHERE id = OLD.content_id;
        END IF;
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_peatch_points_trigger
    AFTER INSERT OR DELETE
    ON likes
    FOR EACH ROW
EXECUTE FUNCTION update_peatch_points();


ALTER TABLE users
    ADD COLUMN likes_count INTEGER DEFAULT 0;
ALTER TABLE posts
    ADD COLUMN likes_count INTEGER DEFAULT 0;
ALTER TABLE collaborations
    ADD COLUMN likes_count INTEGER DEFAULT 0;