create table if not exists settings (
    id int primary key generated always as identity,
    user_id int not null unique references "user"(id) on delete cascade,
    avatar_url text not null default 'avatar/default.png',
    first_name text,
    last_name text,
    updated_at timestamp default null
);

create trigger update_timestamp
    before update on settings
    for each row execute function new_updated_at();

create index if not exists idx_settings_user_id on settings(user_id);

---- create above / drop below ----

drop index if exists idx_settings_user_id;

drop table if exists settings cascade;

drop trigger if exists update_timestamp on settings;
