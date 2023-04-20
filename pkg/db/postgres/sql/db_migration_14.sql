alter table outgoing_webhooks add column secrets varchar null;
update outgoing_webhooks set secrets = config;
