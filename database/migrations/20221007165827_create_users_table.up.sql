create table "users"
(
    "id"                bigserial primary key          not null,
    "first_name"        varchar(255)                   not null,
    "last_name"         varchar(255)                   not null,
    "email"             varchar(255)                   null,
    "email_verified_at" timestamp(0) without time zone null,
    "phone"             varchar(255)                   not null,
    "phone_verified_at" timestamp(0) without time zone null,
    "password"          varchar(255)                   not null,
    "is_activated"      boolean                        not null default true,
    "created_at"        timestamp(0) without time zone null,
    "updated_at"        timestamp(0) without time zone null
);
alter table "users"
    add constraint "users_email_unique" unique ("email");
alter table "users"
    add constraint "users_phone_unique" unique ("phone");