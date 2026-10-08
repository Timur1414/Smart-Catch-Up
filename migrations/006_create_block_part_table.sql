create table if not exists block_part (
    id int primary key generated always as identity,
    block_id int not null references block(id) on delete cascade,
    cluster cluster_type_enum not null,
    start_at timestamptz not null,
    end_at timestamptz not null,

    constraint end_not_before_start check ( end_at >= start_at )
);

---- create above / drop below ----

drop table if exists block_part cascade;
