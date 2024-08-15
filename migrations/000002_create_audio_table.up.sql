CREATE TABLE IF NOT EXISTS audio (
   uuid uuid UNIQUE PRIMARY KEY,
   name VARCHAR (200) NOT NULL,
   description TEXT,
   path VARCHAR (200) NOT NULL
);