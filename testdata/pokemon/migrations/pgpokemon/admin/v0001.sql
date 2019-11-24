-- +migrate Up

CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;

COMMENT ON EXTENSION citext IS 'data type for case-insensitive character strings';

-- +migrate Down

DROP EXTENSION IF EXISTS citext;