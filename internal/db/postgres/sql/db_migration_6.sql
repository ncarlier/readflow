alter table articles add column search_vectors tsvector;
update articles set search_vectors =
  setweight(to_tsvector(substring(coalesce(title, '') for 1000000)), 'A') ||
  setweight(to_tsvector(substring(coalesce(text,  '') for 1000000)), 'B') ||
  setweight(to_tsvector(substring(coalesce(html,  '') for 1000000)), 'C');
create index articles_search_vectors_idx on articles using gin(search_vectors);
