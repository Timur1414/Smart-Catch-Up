create table if not exists notification_action (
    id int primary key generated always as identity,
    notification_id int not null references notification(id) on delete cascade,
    action_type text not null,
    action_target text not null
);

---- create above / drop below ----

drop table if exists notification_action cascade;
