-- SQLite schema with sqlite-vec support
-- Requires sqlite-vec extension to be loaded

-- Cities table (reference data)
CREATE TABLE IF NOT EXISTS cities
(
    id           TEXT PRIMARY KEY,
    name         TEXT NOT NULL,
    country_code TEXT NOT NULL,
    country_name TEXT NOT NULL,
    latitude     REAL,
    longitude    REAL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Opportunities table (reference data)
CREATE TABLE IF NOT EXISTS opportunities
(
    id             TEXT PRIMARY KEY,
    text_en        TEXT NOT NULL,
    text_ru        TEXT NOT NULL,
    description_en TEXT,
    description_ru TEXT,
    icon           TEXT,
    color          TEXT,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Badges table (reference data)
CREATE TABLE IF NOT EXISTS badges
(
    id         TEXT PRIMARY KEY,
    text       TEXT NOT NULL,
    icon       TEXT,
    color      TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Users table
CREATE TABLE IF NOT EXISTS users
(
    id                       TEXT PRIMARY KEY,
    name                     TEXT,
    chat_id                  INTEGER UNIQUE,
    username                 TEXT UNIQUE,
    created_at               TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    notifications_enabled_at TIMESTAMP,
    hidden_at                TIMESTAMP,
    avatar_url               TEXT,
    title                    TEXT,
    description              TEXT,
    language_code            TEXT      DEFAULT 'en',
    last_active_at           TIMESTAMP,
    verification_status      TEXT      DEFAULT 'unverified',
    verified_at              TIMESTAMP,
    embedding_updated_at     TIMESTAMP,
    login_metadata           TEXT,
    location                 TEXT,
    links                    TEXT, -- JSON array of links
    badges                   TEXT, -- JSON array of badge objects
    opportunities            TEXT, -- JSON array of opportunity objects
    CHECK (verification_status IN ('pending', 'verified', 'denied', 'blocked', 'unverified'))
);

-- User vector embeddings using sqlite-vec
CREATE VIRTUAL TABLE IF NOT EXISTS user_embeddings USING vec0
(
    user_id TEXT PRIMARY KEY,
    embedding FLOAT [1536] -- 1536 dimensions for OpenAI embeddings
);

-- Opportunity vector embeddings using sqlite-vec
CREATE VIRTUAL TABLE IF NOT EXISTS opportunity_embeddings USING vec0
(
    opportunity_id TEXT PRIMARY KEY,
    embedding FLOAT [1536]
);


-- User followers (with TTL support)
CREATE TABLE IF NOT EXISTS user_followers
(
    id          TEXT PRIMARY KEY,
    user_id     TEXT      NOT NULL,
    follower_id TEXT      NOT NULL,
    expires_at  TIMESTAMP NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, follower_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (follower_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Collaborations table
CREATE TABLE IF NOT EXISTS collaborations
(
    id                  TEXT PRIMARY KEY,
    user_id             TEXT NOT NULL,
    title               TEXT NOT NULL,
    description         TEXT NOT NULL,
    is_payable          BOOLEAN   DEFAULT FALSE,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    hidden_at           TIMESTAMP,
    opportunity_id      TEXT NOT NULL,
    location            TEXT,
    links               TEXT, -- JSON array of links
    badges              TEXT, -- JSON array of badges
    opportunity         TEXT, -- JSON object with opportunity details
    verification_status TEXT      DEFAULT 'pending',
    verified_at         TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id),
    CHECK (verification_status IN ('pending', 'verified', 'denied', 'blocked', 'unverified'))
);

-- Collaboration interests (with TTL support)
CREATE TABLE IF NOT EXISTS collaboration_interests
(
    id               TEXT PRIMARY KEY,
    user_id          TEXT      NOT NULL,
    collaboration_id TEXT      NOT NULL,
    expires_at       TIMESTAMP NOT NULL,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, collaboration_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (collaboration_id) REFERENCES collaborations (id) ON DELETE CASCADE
);

-- Admins table
CREATE TABLE IF NOT EXISTS admins
(
    id            TEXT PRIMARY KEY,
    username      TEXT UNIQUE NOT NULL,
    password_hash TEXT        NOT NULL,
    chat_id       INTEGER UNIQUE,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Trigger to clean up expired records (can be called periodically)
-- Note: SQLite doesn't have built-in scheduled jobs, this needs to be called externally
CREATE TABLE IF NOT EXISTS cleanup_log
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    table_name  TEXT    NOT NULL,
    deleted     INTEGER NOT NULL,
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);
CREATE INDEX IF NOT EXISTS idx_users_chat_id ON users (chat_id);
CREATE INDEX IF NOT EXISTS idx_users_verification ON users (verification_status, hidden_at);

CREATE INDEX IF NOT EXISTS idx_collaborations_user ON collaborations (user_id);
CREATE INDEX IF NOT EXISTS idx_collaborations_verification ON collaborations (verification_status, hidden_at);
CREATE INDEX IF NOT EXISTS idx_collaborations_created ON collaborations (created_at);

CREATE INDEX IF NOT EXISTS idx_user_followers_expires ON user_followers (expires_at);
CREATE INDEX IF NOT EXISTS idx_collab_interests_expires ON collaboration_interests (expires_at);
