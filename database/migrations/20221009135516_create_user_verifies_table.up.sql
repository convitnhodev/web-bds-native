create table "user_verifies"
(
    "phone_number" varchar(255)                   not null default '',
    "token"        varchar(255)                   not null default '',
    "created_at"   timestamp(0) without time zone not null default '0001-01-01 00:00:00'::timestamp without time zone
);
create index "user_verifies_phone_index" on "user_verifies" ("phone_number");
