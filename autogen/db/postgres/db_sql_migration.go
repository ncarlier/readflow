// Code generated. DO NOT EDIT!

package postgres

// DatabaseSQLMigration is generated form a fileset
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
	"db_migration_3": `alter table users add column plan varchar not null default 'default';
`,
	"db_migration_4": `alter table categories add column rule text;

update categories set rule=rules.rule from rules where rules.category_id = categories.id;
`,
	"db_migration_5": `drop table rules;

alter table articles add column starred boolean not null default 'f';
`,
}

// DatabaseSQLMigrationChecksums is generated from a fileset and contains files checksums
var DatabaseSQLMigrationChecksums = map[string]string{
	"db_migration_1": "6b7ac5c1474bc400c1bbb642fcf3c161f51de7252350eaa261cb1ed796e72b67",
	"db_migration_2": "0be0d1ef1e9481d61db425a7d54378f3667c091949525b9c285b18660b6e8a1d",
	"db_migration_3": "5cd0d3628d990556c0b85739fd376c42244da7e98b66852b6411d27eda20c3fc",
	"db_migration_4": "d5fb83c15b523f15291310ff27d36c099c4ba68de2fd901c5ef5b70a18fedf65",
	"db_migration_5": "16657738407dc4a05c8e2814536078ff598647eb289dfb3aead73f0ac454793b",
}
