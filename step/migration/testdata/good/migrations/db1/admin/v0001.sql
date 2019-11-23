-- +migrate Up
CREATE TABLE admin ( 
    id INT, 
    name VARCHAR(255)
);

-- +migrate Down
DROP TABLE admin;