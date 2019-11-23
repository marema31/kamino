-- +migrate Up
CREATE TABLE user2 ( 
    id INT, 
    name VARCHAR(255)
);

-- +migrate Down
DROP TABLE user2;