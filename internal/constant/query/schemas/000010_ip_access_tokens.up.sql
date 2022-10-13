CREATE TABLE ip_access_tokens
(
    id            uuid PRIMARY KEY     default gen_random_uuid(),
    user_id       uuid        NOT NULL,
    sub_id        string      NOT NULL,
    ip_id         uuid        NOT NULL,
    token         varchar     NOT NULL,
    refresh_token varchar,
    status        VARCHAR(50)          DEFAULT 'ACTIVE',
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT ip_id_fkey FOREIGN KEY (ip_id) REFERENCES identity_providers (id) ON DELETE CASCADE
);