-- +migrate Up
ALTER TABLE pokemon ADD COLUMN defense integer, ADD COLUMN defense_speed integer;

-- +migrate Down

ALTER TABLE pokemon DROP COLUMN defense, DROP COLUMN defense_speed;
