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
    language_code            VARCHAR(2)    NOT NULL DEFAULT 'en',
    country                  VARCHAR(255),
    city                     VARCHAR(255),
    country_code             VARCHAR(2),
    review_status            VARCHAR(255)  NOT NULL DEFAULT 'pending',
    referrer_id              INTEGER REFERENCES users (id),
    last_check_in            timestamptz,
    rating                   INTEGER
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
    country_code   VARCHAR(2)
);

CREATE TABLE user_followers
(
    user_id     INTEGER REFERENCES users (id),
    follower_id INTEGER REFERENCES users (id),
    created_at  timestamptz DEFAULT NOW(),
    UNIQUE (user_id, follower_id)
);

CREATE TABLE collaboration_followers
(
    collaboration_id INTEGER REFERENCES collaborations (id),
    follower_id      INTEGER REFERENCES users (id),
    created_at       timestamptz DEFAULT NOW(),
    UNIQUE (collaboration_id, follower_id)
);

CREATE table collaboration_badges
(
    collaboration_id INTEGER REFERENCES collaborations (id) ON DELETE CASCADE,
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
    hidden_at    timestamptz,
    created_at   timestamptz NOT NULL DEFAULT NOW(),
    updated_at   timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX users_chat_id_index ON users (chat_id);
CREATE UNIQUE INDEX users_username_index ON users (username);
CREATE INDEX review_status_index ON users (review_status);

---
ALTER TABLE users
    DROP COLUMN likes_count,
    DROP COLUMN followers_count,
    DROP COLUMN following_count,
    DROP COLUMN requests_count,
    DROP COLUMN peatch_points;

ALTER TABLE collaborations
    DROP COLUMN likes_count,
    DROP COLUMN requests_count;

DROP TRIGGER IF EXISTS update_referrer_points_trigger ON users;
DROP FUNCTION IF EXISTS update_referrer_points;

DROP TRIGGER IF EXISTS update_follower_counts_trigger ON user_followers;
DROP FUNCTION IF EXISTS update_follower_counts;

DROP TRIGGER IF EXISTS trigger_update_collaborations_hidden_at ON users;
DROP FUNCTION IF EXISTS update_collaborations_hidden_at;

DROP TRIGGER IF EXISTS update_peatch_points_trigger ON likes;
DROP FUNCTION IF EXISTS update_peatch_points;

DROP TABLE IF EXISTS user_interactions;

DROP VIEW IF EXISTS user_activity_feed;
DROP TABLE IF EXISTS collaboration_requests;
DROP TABLE IF EXISTS user_collaboration_requests;
DROP TYPE IF EXISTS collaboration_request_status;


DROP TABLE IF EXISTS likes;
DROP INDEX IF EXISTS likes_user_id_index;
DROP INDEX IF EXISTS likes_content_id_index;
DROP INDEX IF EXISTS likes_content_type_index;

CREATE TABLE collaboration_followers
(
    collaboration_id INTEGER REFERENCES collaborations (id),
    follower_id      INTEGER REFERENCES users (id),
    created_at       timestamptz DEFAULT NOW(),
    UNIQUE (collaboration_id, follower_id)
);

CREATE INDEX idx_users_published_hidden_rating ON users (rating) WHERE published_at IS NOT NULL AND hidden_at IS NULL;
CREATE INDEX idx_users_published_hidden_created_at ON users (created_at DESC) WHERE published_at IS NOT NULL AND hidden_at IS NULL;


alter table public.opportunities
    add text_ru varchar(255);

alter table public.opportunities
    add description_ru text;

