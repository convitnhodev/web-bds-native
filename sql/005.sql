-- chắc kèo phone token và email token ko bao giờ duplicated
create unique index if not exists idx_users_phone_token on users  using btree(phone_token);
create unique index if not exists idx_users_email_token on users  using btree(email_token);