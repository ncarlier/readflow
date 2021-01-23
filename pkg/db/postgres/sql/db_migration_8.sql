alter table articles add column stars int not null default 0;
update articles set stars=1 where starred = true;
alter table articles drop column starred;
