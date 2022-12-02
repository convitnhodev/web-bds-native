alter table products add column if not exists created_by int not null default 0;
CREATE INDEX idx_products_created_by ON products USING btree(created_by);

alter table products add column if not exists is_censorship boolean not null default False;
CREATE INDEX idx_products_is_censorship ON products USING btree(is_censorship);
alter table products add column if not exists censored_at timestamptz;