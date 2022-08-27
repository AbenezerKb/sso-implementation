CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    refresh_token varchar(255) NOT NULL,
    code varchar(255) NOT NULL,
    user_id UUID NOT NULL,
    scope varchar(255),
    redirect_uri varchar(255),
    expires_at timestamptz NOT NULL,
    client_id UUID NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT refreshtoken_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);