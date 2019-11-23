-- +migrate Up
ALTER TABLE pokemon ADD COLUMN attack INT, ADD COLUMN attack_speed INT;

-- +migrate Down

ALTER TABLE pokemon DROP COLUMN attack, DROP COLUMN attack_speed;
