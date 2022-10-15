create table "permissions"
(
    "id"          bigserial primary key,
    "name"        varchar(255) not null default '',
    "slug"        varchar(255) not null default '',
    "description" varchar(255) not null default ''
);
alter table "permissions"
    add constraint "permissions_slug_unique" unique ("slug");
