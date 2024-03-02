alter table incoming_webhooks add column script varchar not null default 'return true;';

alter table categories drop column rule;
alter table categories drop column notification_strategy;
drop type notification_strategy_type;
