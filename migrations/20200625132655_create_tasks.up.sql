CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS groups (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS tasks(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    group_id uuid REFERENCES groups(id)
);

CREATE TABLE IF NOT EXISTS time_frames(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id uuid REFERENCES tasks(id),
    from_time timestamp NOT NULL,
    to_time timestamp NOT NULL
);