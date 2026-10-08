create table if not exists "user" (
    id int primary key generated always as identity,
    is_staff boolean not null default false,
    active boolean not null default true,
    email text not null unique,
    password text not null,
    created_at timestamp not null default now(),
    updated_at timestamp default null,

    constraint password_length check ( length(password) >= 8 ),
    constraint email_is_correct check ( email ~* '^[A-Za-zа-яёА-ЯЁ0-9._%+-]+@[A-Za-zа-яёА-ЯЁ0-9.-]+\.[A-Za-zа-яёА-ЯЁ]{2,}$' ),
    constraint update_at_not_in_past check ( updated_at is null or updated_at >= created_at )
);

create trigger update_timestamp
    before update on "user"
    for each row execute function new_updated_at();

create index if not exists idx_user_email on "user"(email);

---- create above / drop below ----

drop index if exists idx_user_email;

drop table if exists "user" cascade;

drop trigger if exists update_timestamp on "user";
