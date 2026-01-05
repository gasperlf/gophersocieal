ALTER TABLE posts
    ADD COLUMN tags varchar(100) [];

COMMENT ON COLUMN posts.tags IS 'Array of tags associated with the post, each tag can be up to 100 characters long.';

ALTER TABLE posts
    ADD COLUMN updated_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW();
COMMENT ON COLUMN posts.updated_at IS 'Timestamp of the last update to the post.';