alter table invoices add column if not exists total_amount int not null default 0;
alter table payments add column if not exists response_data text not null default '';
