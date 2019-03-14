// Code generated. DO NOT EDIT!

package postgres

var DatabaseSQLMigration = map[string]string{
	"db_migration_1": `create table schema_version (
  version text not null
);

create table users (
  id serial not null,
  username text not null unique,
  last_login_at timestamp with time zone,
  created_at timestamp with time zone not null default now(),

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

create table rules (
  id serial not null,
  user_id int not null,
  priority int not null,
  rule text not null,
  action text not null,
  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone null,

  primary key (id),
  foreign key (user_id) references users(id) on delete cascade
)
`,
}

var DatabaseSQLMigrationChecksums = map[string]string{
	"db_migration_1": "665b35b19bffae3908aa3317e2b98a8bfd46262f0e7e37744cba81833cb356e1",
}
