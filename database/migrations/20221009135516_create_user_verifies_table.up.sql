create table "user_verifies"
(
    "phone_number" varchar(255)                   not null,
    "token"        varchar(255)                   not null,
    "created_at"   timestamp(0) without time zone null
);
create index "user_verifies_phone_index" on "user_verifies" ("phone_number");
