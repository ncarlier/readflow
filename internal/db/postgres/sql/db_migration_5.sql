drop table rules;

alter table articles add column starred boolean not null default 'f';
