alter table "user"
    add column if not exists totp_secret text not null default '',
    add column if not exists totp_enabled boolean not null default false,
    add column if not exists totp_backup_codes text[] not null default '{}';

---- create above / drop below ----

alter table "user"
    drop column if exists totp_backup_codes,
    drop column if exists totp_enabled,
    drop column if exists totp_secret;
