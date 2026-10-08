create table if not exists notification (
    id int primary key generated always as identity,
    cluster cluster_type_enum default null,
    notification_type notification_type_enum not null,
    recipient_id int not null references "user"(id) on delete cascade,
    actor_id int not null references "user"(id) on delete cascade,
    actor_name text not null,
    object_id int not null default 0,
    object_type text not null default '',
    created_at timestamptz not null default now(),
    read_at timestamptz default null,
    payload text not null,

    constraint read_at_not_in_past check ( read_at is null or read_at >= created_at )
);

create index if not exists idx_notification_cluster on notification(cluster);

---- create above / drop below ----

drop index if exists idx_notification_cluster;

drop table if exists notification cascade;
