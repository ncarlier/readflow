create type notification_strategy_type as enum('none', 'individual', 'global');
alter table categories add column notification_strategy notification_strategy_type not null default 'none';
