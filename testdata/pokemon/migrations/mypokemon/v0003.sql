-- +migrate Up
ALTER TABLE pokemon ADD COLUMN defense INT, ADD COLUMN defense_speed INT;

-- +migrate Down

ALTER TABLE pokemon DROP COLUMN defense, DROP COLUMN defense_speed;
