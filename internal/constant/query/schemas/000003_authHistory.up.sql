CREATE TABLE auth_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code varchar(255) NOT NULL,
    user_id UUID NOT NULL,
    scope varchar(255),
    status varchar(255) NOT NULL,
    redirect_uri varchar(255),
    client_id UUID NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT auth_history_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
); 