CREATE TABLE tags (
    id VARCHAR(27) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by_id VARCHAR(27) ,
    edited_by_id VARCHAR(27) NULL,

    name VARCHAR NOT NULL ,
    slug VARCHAR NOT NULL UNIQUE,

    FOREIGN KEY (created_by_id) REFERENCES users(id),
    FOREIGN KEY (edited_by_id) REFERENCES users(id)
);