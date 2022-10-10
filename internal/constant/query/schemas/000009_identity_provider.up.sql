CREATE TABLE identity_providers
(
    id                     uuid PRIMARY KEY     default gen_random_uuid(),
    name                   varchar     NOT NULL,
    logo_url               varchar,
    client_id              varchar     NOT NULL,
    client_secret          varchar     NOT NULL,
    redirect_uri           varchar     NOT NULL,
    authorization_uri      varchar     NOT NULL,
    token_endpoint_url     varchar     NOT NULL,
    user_info_endpoint_url varchar,
    status                 VARCHAR(255)         DEFAULT 'ACTIVE',
    created_at             timestamptz NOT NULL DEFAULT now(),
    updated_at             timestamptz NOT NULL DEFAULT now()
);
