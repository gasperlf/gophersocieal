CREATE TABLE user_invitations (
    token bytea PRIMARY KEY,
    user_id bigint NOT NULL
);

