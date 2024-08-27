CREATE TABLE IF NOT EXISTS audio (
   id uuid UNIQUE PRIMARY KEY,
   name VARCHAR (200) NOT NULL,
   description TEXT,
   path VARCHAR (200) NOT NULL
   owner uuid NULL
);