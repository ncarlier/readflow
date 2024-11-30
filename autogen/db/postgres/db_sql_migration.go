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
	"db_migration_10": `ALTER TYPE article_status RENAME VALUE 'unread' TO 'inbox';
ALTER TYPE article_status ADD VALUE 'to_read' AFTER 'read';
`,
	"db_migration_11": `create type notification_strategy_type as enum('none', 'individual', 'global');
alter table categories add column notification_strategy notification_strategy_type not null default 'none';
`,
	"db_migration_12": `alter table devices add column last_seen_at timestamp with time zone not null default now();
`,
	"db_migration_13": `alter table incoming_webhooks add column script varchar not null default 'return true;';

alter table categories drop column rule;
alter table categories drop column notification_strategy;
drop type notification_strategy_type;
`,
	"db_migration_14": `alter table outgoing_webhooks add column secrets varchar null;
update outgoing_webhooks set secrets = config;
`,
	"db_migration_15": `alter table articles add column thumbhash varchar null;
`,
	"db_migration_16": `alter table articles add column tags text[] not null default '{}';
create index tags_index on articles using gin (tags);

alter table articles add column search tsvector generated always as (
  setweight(array_to_tsvector(tags), 'A') || ' ' ||
  setweight(to_tsvector('english', title), 'B') || ' ' ||
  setweight(to_tsvector('english', coalesce(text, '')), 'C') || ' ' ||
  setweight(to_tsvector('english', substring(coalesce(html, '') for 1000000)), 'D')
) STORED;

create index articles_search_idx on articles using gin(search);

drop index articles_search_vectors_idx;
alter table articles drop column search_vectors;
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
	"db_migration_6": `alter table articles add column search_vectors tsvector;
update articles set search_vectors =
  setweight(to_tsvector(substring(coalesce(title, '') for 1000000)), 'A') ||
  setweight(to_tsvector(substring(coalesce(text,  '') for 1000000)), 'B') ||
  setweight(to_tsvector(substring(coalesce(html,  '') for 1000000)), 'C');
create index articles_search_vectors_idx on articles using gin(search_vectors);
`,
	"db_migration_7": `alter table archivers rename to outgoing_webhooks;
alter table api_keys rename to incoming_webhooks;
update outgoing_webhooks set provider='generic' where provider = 'webhook';
`,
	"db_migration_8": `alter table articles add column stars int not null default 0;
update articles set stars=1 where starred = true;
alter table articles drop column starred;
`,
	"db_migration_9": `alter table users add column customer_id varchar not null default '';
`,
}

// DatabaseSQLMigrationChecksums is generated from a fileset and contains files checksums
var DatabaseSQLMigrationChecksums = map[string]string{
	"db_migration_1":  "6b7ac5c1474bc400c1bbb642fcf3c161f51de7252350eaa261cb1ed796e72b67",
	"db_migration_10": "935f7f7208d0230865d0915bf8f6b940331084d3aeb951536605f879a85a842f",
	"db_migration_11": "1150b8fa81099eb5956989560e8eebecafe5e39cbd1a5f6f7d23f3dfceb810bf",
	"db_migration_12": "b24497bb03f04fb4705ae752f8a5bf69dad26f168bc8ec196af93aee29deef49",
	"db_migration_13": "4a52465eeb50a236d7f7a94cc51cd78238de0f885a6d29da4a548b5c389ebe81",
	"db_migration_14": "f2c6e03988386e662f943d0f37255cf6db19b69e2c4f63c312f3778b401bb96a",
	"db_migration_15": "edf9f683832d4b5c8c0d681f479750794ca19aea115a89b69700d4f415104fc3",
	"db_migration_16": "6a97e3bd2e1b238fb40524f095810703f7c4a873f43723f350b06661eae995b7",
	"db_migration_2":  "0be0d1ef1e9481d61db425a7d54378f3667c091949525b9c285b18660b6e8a1d",
	"db_migration_3":  "5cd0d3628d990556c0b85739fd376c42244da7e98b66852b6411d27eda20c3fc",
	"db_migration_4":  "d5fb83c15b523f15291310ff27d36c099c4ba68de2fd901c5ef5b70a18fedf65",
	"db_migration_5":  "16657738407dc4a05c8e2814536078ff598647eb289dfb3aead73f0ac454793b",
	"db_migration_6":  "82606f963d687906ec932d2a6021a29b0d1480260c8a1f7fe7da8edfad8bfbf5",
	"db_migration_7":  "05329d34279e8787592c48e97164dd0be0a1f42835da3f4aa129819296828a8d",
	"db_migration_8":  "36dfebfaec092e686472a440a7d22e318a4f46567d18864ec3e53b94ac12e837",
	"db_migration_9":  "ecfa9532599414b3e79dad336e7069b22f888b981db90b067104026ecb0a56ac",
}
