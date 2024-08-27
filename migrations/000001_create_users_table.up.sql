CREATE TABLE IF NOT EXISTS users(
   id uuid UNIQUE PRIMARY KEY,

   tg_id BIGINT NOT NULL,
   username VARCHAR (100) NOT NULL,
   first_name VARCHAR (100) NULL,
   last_name VARCHAR (100) NULL,
   photo_url VARCHAR (255) NULL,

   provider VARCHAR (50) NOT NULL,
   role VARCHAR (10)
);