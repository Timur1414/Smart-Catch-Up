create table if not exists block_part (
    id int primary key generated always as identity,
    block_id int not null references block(id) on delete cascade,
    cluster text not null,
    text text not null,
    start_at timestamptz not null,
    end_at timestamptz not null,

    constraint end_not_before_start check ( end_at >= start_at )
);

create index if not exists idx_block_part_block_id on block_part (block_id);

---- create above / drop below ----

drop index if exists idx_block_part_block_id;

drop table if exists block_part cascade;
