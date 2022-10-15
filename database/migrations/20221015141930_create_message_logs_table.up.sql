create table "message_logs"
(
    "id"           bigserial primary key,
    "phone_number" varchar(255)                   not null default '',
    "content"      varchar(255)                   not null default '',
    "request_id"   varchar(255)                   not null default '',
    "sms_id"       varchar(255)                   not null default '',
    "code_result"  varchar(255)                   not null default '',
    "sent_status"  varchar(255)                   not null default '',
    "telco_id"     varchar(255)                   not null default '',
    "created_at"   timestamp(0) without time zone not null default '0001-01-01 00:00:00'::timestamp without time zone,
    "updated_at"   timestamp(0) without time zone not null default '0001-01-01 00:00:00'::timestamp without time zone
);
alter table "message_logs"
    add constraint "message_logs_request_id_unique" unique ("request_id");
alter table "message_logs"
    add constraint "message_logs_sms_id_unique" unique ("sms_id");
