CREATE TABLE sessions (
  id SERIAL PRIMARY KEY,  -- Auto-incrementing integer for unique ID
  username VARCHAR(255) NOT NULL UNIQUE,
  refresh_token VARCHAR NOT NULL 
);