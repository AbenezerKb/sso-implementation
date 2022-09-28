CREATE TABLE resource_servers
(
    id         UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    name       varchar(255) NOT NULL UNIQUE,
    created_at timestamptz  NOT NULL DEFAULT now(),
    updated_at timestamptz  NOT NULL DEFAULT now()
);