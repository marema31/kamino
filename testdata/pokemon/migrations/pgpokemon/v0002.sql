-- +migrate Up
ALTER TABLE pokemon ADD COLUMN attack integer, ADD COLUMN attack_speed integer;

-- +migrate Down

ALTER TABLE pokemon DROP COLUMN attack, DROP COLUMN attack_speed;
