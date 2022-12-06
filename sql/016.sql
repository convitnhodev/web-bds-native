alter table products add column if not exists remain_of_slot int not null default 0;
alter table products drop column if exists escrow_amount;
alter table products add column if not exists deposit_percent decimal(4, 2) not null default 0;