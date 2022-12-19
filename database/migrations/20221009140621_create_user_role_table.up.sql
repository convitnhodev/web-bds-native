create table "users_roles"
(
    "user_id" bigint not null default 0,
    "role_id" bigint not null default 0
);
alter table "users_roles"
    add constraint "users_roles_user_id_role_id_unique" unique ("user_id", "role_id");
