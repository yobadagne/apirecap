CREATE TABLE users (
  id SERIAL PRIMARY KEY,  -- Auto-incrementing integer for unique ID
  username VARCHAR(255) NOT NULL UNIQUE, -- Username with max length of 255, must be unique
  password VARCHAR(255) NOT NULL, -- Password field
  email VARCHAR(255) NOT NULL UNIQUE -- Email address, must be unique
);