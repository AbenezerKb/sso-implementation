CREATE TABLE authcodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code varchar(255) NOT NULL,
    user_id UUID NOT NULL,
    status varchar(255) NOT NULL DEFAULT 'new',
    scope varchar(255) NOT NULL DEFAULT 'email',
    redirect_uri varchar(255) NOT NULL DEFAULT '',
    client_id UUID NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT authcode_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
); 