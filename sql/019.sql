

alter table invoices add column if not exists slot_canceled_by int not null default 0;

ALTER TYPE invoice_status RENAME VALUE 'canceled' to 'slot_canceled';
ALTER TYPE invoice_status RENAME VALUE 'paid' to 'collect_completed';
ALTER TYPE invoice_status ADD VALUE 'collecting';
ALTER TYPE invoice_status ADD VALUE 'collect_canceled';

ALTER TYPE payment_method RENAME VALUE 'appotapay' to 'appotapay_payment';
ALTER TYPE payment_method RENAME VALUE 'bank_transfer' to 'appotapay_bill';

alter table payments add column if not exists appotapay_trans_id text not null default '';
alter table payments add column if not exists refund_id text not null default '';
alter table payments add column if not exists refund_response text not null default '';
alter table payments add column if not exists transaction_at timestamptz;
alter table payments add column if not exists refund_at timestamptz;


alter table payments add column if not exists appotapay_bill_code text not null default '';
alter table payments add column if not exists appotapay_account_no text not null default '';
alter table payments add column if not exists appotapay_account_name text not null default '';
alter table payments add column if not exists appotapay_bank_code text not null default '';
alter table payments add column if not exists appotapay_bank_name text not null default '';
alter table payments add column if not exists appotapay_bank_branch text not null default '';
