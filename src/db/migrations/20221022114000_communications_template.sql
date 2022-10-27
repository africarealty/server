-- +goose Up
-- +goose StatementBegin
set schema 'africarealty';

create table templates
(
    id varchar primary key,
    title varchar not null,
    body  text    not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp
);
insert into templates (id, title, body) values ('test','Title push test','Some text with placeholder {{TestName}}', '2021-01-01 00:00:00', '2021-01-01 00:00:00', NULL );
insert into templates (id, title, body) values ('auth.registration-activation', 'AfricaRealty registration', 'Hi {{Name}},\n\nClick the link below to confirm your registration process!\n\n{{RegistrationLink}}\n\nThanks,\n\nAfrica Realty', '2021-01-01 00:00:00', '2021-01-01 00:00:00', NULL );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
set schema 'africarealty';

drop table templates;
-- +goose StatementEnd
