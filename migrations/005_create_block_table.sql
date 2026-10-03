create table if not exists block (
    id int primary key generated always as identity,
    user_id int not null references "user"(id) on delete cascade,
    start_at timestamp not null,
    end_at timestamp not null,
    status block_status_enum not null default 'created',

    constraint end_not_before_start check ( end_at >= start_at )
);

---- create above / drop below ----

drop table if exists block cascade;
