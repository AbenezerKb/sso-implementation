CREATE TABLE internalrefreshtokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    refreshtoken varchar(255) NOT NULL,
    user_id UUID NOT NULL,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT internal_refreshtoken_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);