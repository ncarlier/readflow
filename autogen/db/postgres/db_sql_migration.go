// Code generated. DO NOT EDIT!

package postgres

var DatabaseSQLMigration = map[string]string{
	"db_migration_1": `create table schema_version (
  version text not null
);


create table users (
  id serial not null,
  username text not null unique,
  enabled bool not null default 't',
  last_login_at timestamp with time zone,
  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone null,

  primary key (id)
);

create table categories (
  id serial not null,
  user_id int not null,
  title text not null,
  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone null,

  primary key (id),
  unique (user_id, title),
  foreign key (user_id) references users(id) on delete cascade
);

create type article_status as enum('unread', 'read');

create table articles (
  id bigserial not null,
  user_id int not null,
  category_id int null,
  title text not null,
  text text,
  html text,
  url text,
  image text,
  hash text not null,
  status article_status default 'unread',
  published_at timestamp with time zone null,
  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone null,
    
  primary key (id),
  unique (user_id, hash),
  foreign key (user_id) references users(id) on delete cascade,
  foreign key (category_id) references categories(id) on delete set null
);

create index articles_user_status_idx on articles(user_id, status);
create index articles_user_status_category_idx on articles(user_id, status, category_id);

create table api_keys (
  id serial not null,
  user_id int not null,
  alias text not null,
  token text not null unique,
  last_usage_at timestamp with time zone,
  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone null,

  primary key (id),
  unique (user_id, alias),
  foreign key (user_id) references users(id) on delete cascade
);

create table archivers (
  id serial not null,
  user_id int not null,
  alias text not null,
  is_default bool not null default 'f',
  provider text not null,
  config json not null,
  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone null,

  primary key (id),
  unique (user_id, alias),
  foreign key (user_id) references users(id) on delete cascade
);

create table rules (
  id serial not null,
  user_id int not null,
  category_id int null,
  alias text not null,
  priority int not null,
  rule text not null,
  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone null,

  primary key (id),
  unique (user_id, alias),
  foreign key (user_id) references users(id) on delete cascade,
  foreign key (category_id) references categories(id) on delete set null
)
`,
	"db_migration_2": `create table devices (
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
`,
}

var DatabaseSQLMigrationChecksums = map[string]string{
	"db_migration_1": "6b7ac5c1474bc400c1bbb642fcf3c161f51de7252350eaa261cb1ed796e72b67",
	"db_migration_2": "0be0d1ef1e9481d61db425a7d54378f3667c091949525b9c285b18660b6e8a1d",
}
