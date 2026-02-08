create table if not exists roles (
    id bigserial primary key,
    name varchar(255) not null unique,
    level int not null default 0 ,
    description text
);


insert into roles (name, level, description) values
('user', 1, 'Create post and comments'),
('moderator', 2, 'A moderator can update other users posts'),
('admin', 3, 'Admin can update and delete other user posts');