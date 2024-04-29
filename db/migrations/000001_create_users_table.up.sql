CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  email VARCHAR(255) NOT NULL UNIQUE,
  last_login TIMESTAMP
);

CREATE TABLE items (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title VARCHAR(255) NOT NULL,
  initial_price DECIMAL(10,2) NOT NULL,
  sold_price DECIMAL(10,2),
  description TEXT
);