alter table block_part add column if not exists rank int not null default 10 check ( rank > 0 );

---- create above / drop below ----

alter table block_part drop column if exists rank;
