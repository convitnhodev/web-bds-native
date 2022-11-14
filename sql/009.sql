
create type file_status as enum ('new', 'sync', 'deleting_local', 'deleted_local', 'deteling_both', 'deleted_both');

create table files (
    id int generated by default as identity primary key,
    local_path text not null default '',
    cloud_link text not null default '',
    status file_status default 'new',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

CREATE INDEX idx_files_local_path ON files USING btree(local_path);
CREATE INDEX idx_files_created_at ON files USING btree(created_at);

create type user_kyc_status as enum ('no_kyc', 'submited_kyc', 'approved_kyc', 'rejected_kyc');
create table kyc (
    id int generated by default as identity primary key,
    user_id int not null,
    front_identity_card text not null default '',
    back_identity_card text not null default '',
    selfie_image text not null default '',
    feedback text not null default '',
    approved_by int,
    rejected_by int,
    status user_kyc_status not null default 'submited_kyc',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

CREATE INDEX idx_kyc_user_id ON kyc USING btree(user_id);

alter table users add column if not exists last_kyc_status user_kyc_status not null default 'no_kyc';
CREATE INDEX idx_users_last_kyc_status ON users USING btree(last_kyc_status);