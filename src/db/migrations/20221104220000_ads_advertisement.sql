-- +goose Up
-- +goose StatementBegin
set schema 'africarealty';

create table advertisements
(
    id           uuid primary key,
    code         varchar   not null,
    user_id      uuid      not null,
    status       varchar   not null,
    sub_status   varchar   not null,
    type         varchar,
    sub_type     varchar   not null,
    details      jsonb,
    activated_at timestamp not null,
    closed_at    timestamp not null,
    created_at   timestamp not null,
    updated_at   timestamp not null,
    deleted_at   timestamp null
);

create index idx_ads_code on advertisements(code);
create index idx_ads_user on advertisements(user_id);

create sequence seq_ads_code start 1;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
set schema 'africarealty';

drop table advertisements;
drop sequence seq_ads_code;
-- +goose StatementEnd