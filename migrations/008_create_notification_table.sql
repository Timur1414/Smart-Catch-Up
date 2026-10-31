create table if not exists notification (
    id int primary key generated always as identity,
    cluster text default null,
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

create index if not exists idx_notification_unclustered_range on notification(recipient_id, created_at) where cluster is null;

---- create above / drop below ----

drop index if exists idx_notification_unclustered_range;

drop index if exists idx_notification_cluster;

drop table if exists notification cascade;
