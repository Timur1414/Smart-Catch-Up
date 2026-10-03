create table if not exists notification_action (
    id int primary key generated always as identity,
    notification_id int not null references notification(id) on delete cascade,
    action_type text not null,
    action_target text not null
);

create index if not exists idx_notification_action_notification_id on notification_action(notification_id);

---- create above / drop below ----

drop index if exists idx_notification_action_notification_id;

drop table if exists notification_action cascade;
