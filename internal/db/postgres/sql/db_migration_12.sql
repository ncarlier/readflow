alter table devices add column last_seen_at timestamp with time zone not null default now();
