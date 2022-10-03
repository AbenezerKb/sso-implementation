CREATE TABLE scopes(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(255) NOT NULL,
    description varchar(255) NOT NULL,
    resource_server_id UUID,
    resource_server_name varchar(255),
    status  varchar(255) NOT NULL DEFAULT 'ACTIVE',
    created_at timestamptz NOT NULL DEFAULT now()
);