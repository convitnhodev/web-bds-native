alter table products add column if not exists slug text not null default '';
alter table products add column if not exists full_content text not null default '';
alter table products add column if not exists product_type text not null default '';
alter table products add column if not exists business_advantage text not null default '';
alter table products add column if not exists legal text not null default '';

alter table products add column if not exists area int not null default 0;