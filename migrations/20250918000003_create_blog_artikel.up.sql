CREATE TABLE blog_artikels (
    id VARCHAR(27) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    title VARCHAR NOT NULL,
    slug VARCHAR UNIQUE NOT NULL,
    content TEXT NOT NULL,
    excerpt TEXT,
    image_url VARCHAR,
    category_id VARCHAR(27) NOT NULL,
    author_id VARCHAR(27) NOT NULL,
    status VARCHAR NOT NULL DEFAULT 'draft',
    source TEXT NOT NULL DEFAULT '-',   -- url path or "-" as default
    featured INT UNIQUE,    
    views BIGINT DEFAULT 0,
    published_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (category_id) REFERENCES categories(id),
    FOREIGN KEY (author_id) REFERENCES users(id)
);

CREATE TABLE article_tags (
    blog_article_id VARCHAR(27) REFERENCES blog_artikels(id) ON DELETE CASCADE,
    tag_id VARCHAR(27) REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (blog_article_id, tag_id)
);