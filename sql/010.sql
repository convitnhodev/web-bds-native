alter table users add column if not exists reset_pwd_token text not null default '';
alter table users add column if not exists rpt_expired_at timestamp with time zone;
