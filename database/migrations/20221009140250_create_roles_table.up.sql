create table "roles"
(
    "id"          bigserial primary key          not null,
    "name"        varchar(255)                   not null,
    "slug"        varchar(255)                   not null,
    "description" varchar(255)                   not null,
    "created_at"  timestamp(0) without time zone null,
    "updated_at"  timestamp(0) without time zone null
);
alter table "roles"
    add constraint "roles_slug_unique" unique ("slug");
