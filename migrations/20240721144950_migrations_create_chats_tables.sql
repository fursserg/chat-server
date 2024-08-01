-- +goose Up
create table chats (
  id serial primary key,
  title text not null,
  user_ids json not null,
  status int not null,
  created_at timestamp not null default now(),
  updated_at timestamp
);

-- +goose Down
drop table chats;
