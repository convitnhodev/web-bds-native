create table "user_role"
(
    "user_id"    bigint                         not null,
    "role_id"    bigint                         not null,
    "created_at" timestamp(0) without time zone null,
    "updated_at" timestamp(0) without time zone null
);
alter table "user_role"
    add constraint "user_role_user_id_role_id_unique" unique ("user_id", "role_id");