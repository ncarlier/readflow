alter table categories add column rule text;

update categories set rule=rules.rule from rules where rules.category_id = categories.id;
