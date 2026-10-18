alter type cluster_type_enum add value if not exists 'important';

alter table settings
    add column if not exists important text[] not null default '{}';

---- create above / drop below ----

alter table settings
    drop column if exists important;
