-- +goose up
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    title VARCHAR(4000) NOT NULL,
    url VARCHAR(4000) NOT NULL,
    description TEXT NULL,
    published_at VARCHAR(100) NULL,
    feed_id UUID NOT NULL,
    UNIQUE (url),
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose down
DROP TABLE posts;
