alter table archivers rename to outgoing_webhooks;
alter table api_keys rename to incoming_webhooks;
update outgoing_webhooks set provider='generic' where provider = 'webhook';
