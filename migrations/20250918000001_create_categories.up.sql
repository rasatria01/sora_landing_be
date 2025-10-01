CREATE TABLE categories (
    id VARCHAR(27) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    name VARCHAR NOT NULL UNIQUE,
    slug VARCHAR NOT NULL UNIQUE,
    created_by_id VARCHAR(27) ,
    edited_by_id VARCHAR(27) NULL,
    FOREIGN KEY (created_by_id) REFERENCES users(id),
    FOREIGN KEY (edited_by_id) REFERENCES users(id)

);