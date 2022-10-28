create unique index if not exists idx_users_email on users  using btree(email) where email <> '';
create unique index if not exists idx_users_phone on users  using btree(phone) where phone <> '';
