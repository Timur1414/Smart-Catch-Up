create table if not exists digest (
    id int primary key generated always as identity,
    user_id int not null unique references "user"(id) on delete cascade,
    block_id int not null unique references block(id) on delete cascade,
    created_at timestamp not null default now(),
    updated_at timestamp,

    constraint update_at_not_in_past check ( updated_at >= created_at )
);

create trigger update_timestamp
    before update on digest
    for each row execute function new_updated_at();

---- create above / drop below ----

drop table if exists digest cascade;

drop trigger if exists update_timestamp on digest;
