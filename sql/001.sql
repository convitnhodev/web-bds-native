create table users (
    id int generated by default as identity primary key,
    email text not null default '',
    phone text not null default '',
    password text not null default '',
    first_name text not null default '',
    last_name text not null default '',
    roles text[] not null default '{}',
    email_token text not null default '',
    phone_token text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default '0001-01-01 00:00:00',
    send_verified_email_at timestamptz not null default '0001-01-01 00:00:00',
    send_verified_phone_at timestamptz not null default '0001-01-01 00:00:00'
)