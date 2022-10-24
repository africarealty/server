-- +goose Up
-- +goose StatementBegin
set schema 'africarealty';

create table emails
(
    id uuid primary key,
    user_id uuid,
    email varchar not null,
    template_id varchar,
    template_data jsonb,
    send_status varchar not null,
    text text not null,
    attachments jsonb,
    error_desc varchar,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp null
);

create index idx_emails_created on emails(created_at);
create index idx_emails_usr on emails(user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
set schema 'africarealty';
drop table emails;
-- +goose StatementEnd