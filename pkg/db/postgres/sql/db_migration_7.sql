alter table archivers rename to outbound_services;
alter table api_keys rename to inbound_services;
alter table inbound_services add column provider text not null default 'webhook';
alter table inbound_services add column config json not null default '{}'::JSON;
