CREATE TABLE IF NOT EXISTS comments (
    id bigserial primary key,
    post_id bigint not null,
    user_id bigint not null,
    content text not null,
    created_at timestamp(0) with time zone not null default now()
);