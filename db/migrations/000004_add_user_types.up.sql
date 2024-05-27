CREATE TABLE user_types (
    id SERIAL PRIMARY KEY,
    type_name VARCHAR(255) NOT NULL,
    type_description TEXT,
    type_code VARCHAR(255) NOT NULL
);

INSERT INTO user_types (type_name, type_description, type_code) VALUES ('Seller', 'Seller of items', 'SELLER');
INSERT INTO user_types (type_name, type_description, type_code) VALUES ('Buyer', 'Buyer of items', 'BUYER');

ALTER TABLE users ADD COLUMN user_type_id INTEGER REFERENCES user_types(id);