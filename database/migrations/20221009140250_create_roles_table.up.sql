create table "roles"
(
    "id"          bigserial primary key,
    "name"        varchar(255) not null default '',
    "slug"        varchar(255) not null default '',
    "description" varchar(255) not null default ''
);
alter table "roles"
    add constraint "roles_slug_unique" unique ("slug");
