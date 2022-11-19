
create type partner_status as enum ('not_yet', 'apply', 'approved', 'rejected');

create table partner (
    id int generated by default as identity primary key,
    user_id int not null,
    message text not null default '',
    cv_link text not null default '',
    status partner_status not null default 'apply',
    feedback text not null default '',
    approved_by int,
    rejected_by int,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

alter table users add column if not exists partner_status partner_status not null default 'not_yet';
CREATE INDEX idx_users_partner_status ON users USING btree(partner_status);