-- +goose Up
CREATE ROLE africarealty LOGIN PASSWORD 'africarealty' NOINHERIT CREATEDB;
CREATE SCHEMA africarealty AUTHORIZATION africarealty;
GRANT USAGE ON SCHEMA africarealty TO PUBLIC;

-- +goose Down
DROP SCHEMA africarealty;
DROP ROLE africarealty;
