create table if not exists block (
    id int primary key generated always as identity,
    user_id int not null references "user"(id) on delete cascade,
    start_at timestamptz not null,
    end_at timestamptz not null,
    status block_status_enum not null default 'created',

    constraint end_not_before_start check ( end_at >= start_at )
);

create index if not exists idx_block_user_id on block(user_id);

---- create above / drop below ----

drop index if exists idx_block_user_id;

drop table if exists block cascade;
