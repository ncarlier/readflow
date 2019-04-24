create table devices (
  id serial not null,
  user_id int not null,
  key text not null,
  subscription json not null,
  created_at timestamp with time zone not null default now(),

  primary key (id),
  unique (user_id, key),
  foreign key (user_id) references users(id) on delete cascade
);

create table properties (
  rev serial not null,
  vapid_public_key text not null,
  vapid_private_key text not null,
  created_at timestamp with time zone not null default now(),

  primary key (rev)
)
