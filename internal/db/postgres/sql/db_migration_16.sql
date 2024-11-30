alter table articles add column tags text[] not null default '{}';
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
