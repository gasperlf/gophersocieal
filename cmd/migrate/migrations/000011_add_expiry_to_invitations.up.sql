alter table user_invitations
add column expiry timestamp(0) with time zone not null;