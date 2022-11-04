-- có thể xóa product
alter table products add column if not exists is_deleted bool not null default false;


alter table attachments add column if not exists product_id int not null default 0;
alter table attachments add column if not exists created_at timestamptz not null default now();
alter table attachments add column if not exists updated_at timestamptz not null default '0001-01-01 00:00:00';
