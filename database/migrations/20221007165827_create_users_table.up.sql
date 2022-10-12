create table "users"
(
    "id"           bigserial primary key          not null,
    "first_name"   varchar(255)                   not null,
    "last_name"    varchar(255)                   not null,
    "email"        varchar(255)                   null,
    "phone_number" varchar(255)                   not null,
    "verified_at"  timestamp(0) without time zone null,
    "password"     varchar(255)                   not null,
    "is_activated" boolean                        not null default true,
    "created_at"   timestamp(0) without time zone null,
    "updated_at"   timestamp(0) without time zone null
);
create index "users_phone_unique" on "users" ("phone_number");
