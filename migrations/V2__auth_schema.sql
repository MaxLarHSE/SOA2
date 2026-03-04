CREATE TABLE users (
                       id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       email         VARCHAR(255) NOT NULL UNIQUE,
                       password_hash VARCHAR(255) NOT NULL,
                       role          VARCHAR(50)  NOT NULL DEFAULT 'USER',
                       created_at    TIMESTAMP    NOT NULL DEFAULT NOW(),
                       updated_at    TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE refresh_tokens (
                                id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                user_id     UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                token_hash  VARCHAR(64)  NOT NULL UNIQUE,
                                expires_at  TIMESTAMP    NOT NULL,
                                revoked     BOOLEAN      NOT NULL DEFAULT false,
                                created_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_user_id    ON refresh_tokens(user_id);