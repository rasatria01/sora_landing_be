CREATE TABLE social_media (
    id VARCHAR(27) PRIMARY KEY,
    name VARCHAR NOT NULL,
    user_name VARCHAR NOT NULL,
    visible BOOLEAN NOT NULL DEFAULT true
);