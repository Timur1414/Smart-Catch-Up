create extension if not exists pg_stat_statements;
create extension if not exists pg_trgm;

---- create above / drop below ----

drop extension if exists pg_trgm;
drop extension if exists pg_stat_statements;
