create table "roles_permissions"
(
    "role_id"       bigint not null default 0,
    "permission_id" bigint not null default 0
);
alter table "roles_permissions"
    add constraint "roles_permissions_role_id_permission_id_unique" unique ("role_id", "permission_id");
