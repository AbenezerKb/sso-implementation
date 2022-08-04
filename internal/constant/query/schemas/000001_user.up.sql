CREATE TABLE users (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "first_name" varchar      NOT NULL,
    "middle_name" varchar      NOT NULL,
    "last_name" varchar      NOT NULL,
    "email" varchar UNIQUE,
    "phone" varchar UNIQUE NOT NULL,
    "password" varchar NOT NULL,
    "user_name" varchar NOT NULL,
    "gender" varchar NOT NULL,
    "profile_picture" varchar,
    "status" varchar DEFAULT 'inactive',
    "created_at" timestamptz NOT NULL
);