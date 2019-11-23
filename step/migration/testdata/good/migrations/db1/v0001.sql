-- +migrate Up
CREATE TABLE user ( 
    id INT, 
    name VARCHAR(255)
);

-- +migrate Down
DROP TABLE user;