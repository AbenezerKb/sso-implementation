CREATE TABLE clients (
                         "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                         "name" varchar NOT NULL,
                         "client_type" varchar NOT NULL,
                         "redirect_uris" string NOT NULL,
                         "scopes" string NOT NULL,
                         "secret" varchar NOT NULL,
                         "logo_url" varchar NOT NULL,
                         "status" varchar NOT NULL DEFAULT 'ACTIVE',
                         "created_at" timestamptz NOT NULL default now()
);