create index if not exists idx_settings_firstname_search on settings using gin(first_name gin_trgm_ops);
create index if not exists idx_settings_lastname_search on settings using gin(last_name gin_trgm_ops);

---- create above / drop below ----

drop index if exists idx_settings_lastname_search;
drop index if exists idx_settings_firstname_search;
