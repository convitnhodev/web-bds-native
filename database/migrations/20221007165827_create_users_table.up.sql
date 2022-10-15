create table "users"
(
    "id"           bigserial primary key,
    "first_name"   varchar(255)                   not null default '',
    "last_name"    varchar(255)                   not null default '',
    "email"        varchar(255)                   not null default '',
    "phone_number" varchar(255)                   not null default '',
    "verified_at"  timestamp(0)                   not null default '0001-01-01 00:00:00'::timestamp without time zone,
    "password"     varchar(255)                   not null default '',
    "is_activated" boolean                        not null default true,
    "created_at"   timestamp(0) without time zone not null default '0001-01-01 00:00:00'::timestamp without time zone,
    "updated_at"   timestamp(0) without time zone not null default '0001-01-01 00:00:00'::timestamp without time zone
);
create unique index "users_phone_unique" on "users" ("phone_number");
create unique index "users_email_unique" on "users" ("email") where "email" <> '';
