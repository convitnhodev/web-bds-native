


alter table products add column if not exists num_of_slot int not null default 0;
alter table products add column if not exists cost_per_slot int not null default 0;
alter table products add column if not exists escrow_amount int not null default 0;


create type invoice_status as enum ('open', 'partial', 'paid', 'void');
create table invoices (
    id int generated by default as identity primary key,
    user_id int not null,
    status invoice_status not null default 'open',
    invoice_synced_at timestamptz,
    invoice_serect text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);
create index idx_invoices_user_id on invoices using btree(user_id);

create table invoice_items (
    id int generated by default as identity primary key,
    invoice_id int not null,
    product_id int not null,
    quatity int not null default 0,
    cost_per_slot int not null default 0,
    amount int not null default 0,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);
create index idx_invoice_items_invoice_id on invoice_items using btree(invoice_id);
create index idx_invoice_items_product_id on invoice_items using btree(product_id);

create type payment_menthod as enum ('appotapay', 'bank_transfer');
create type payment_type as enum ('escrow', 'full');
create type transaction_type as enum ('pay', 'refund');

create table payments (
    id int generated by default as identity primary key,
    invoice_id int not null,
    amount int not null default 0,
    menthod payment_menthod not null default 'appotapay',
    pay_type payment_type not null default 'full',
    tx_type transaction_type not null default 'pay',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index idx_payments_invoice_id on payments using btree(invoice_id);