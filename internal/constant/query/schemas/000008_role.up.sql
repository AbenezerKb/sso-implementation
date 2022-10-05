CREATE TABLE roles
(
    name       VARCHAR(255) PRIMARY KEY,
    status     VARCHAR(255) DEFAULT 'ACTIVE',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);