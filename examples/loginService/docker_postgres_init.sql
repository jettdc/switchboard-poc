CREATE TABLE person
(
    username VARCHAR(50) NOT NULL PRIMARY KEY, 
    password VARCHAR(50) NOT NULL, 
    token VARCHAR(500) NOT NULL
);

INSERT INTO person(username, password, token) VALUES
 ('yich7110', '1234', 'foo_bar');