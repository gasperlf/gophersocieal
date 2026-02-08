alter table if exists users
add column if not exists role_id bigint references roles(id) default 1;  


update users set role_id = (select id from roles where name = 'user');

alter table users
alter column role_id drop default;

alter table users
alter column role_id set not null;