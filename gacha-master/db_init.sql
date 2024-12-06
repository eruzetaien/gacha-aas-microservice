CREATE DATABASE gacha_master;

CREATE TABLE users (
   id SERIAL PRIMARY KEY,
   name VARCHAR(100) NOT NULL UNIQUE,
   username VARCHAR(100) NOT NULL UNIQUE,
   password VARCHAR(60) NOT NULL,
   created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE gacha_system (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  name VARCHAR(100) NOT NULL,
  endpoint TEXT NOT NULL UNIQUE,
  endpoint_id TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
  FOREIGN KEY (user_id)
      REFERENCES users(id)
      ON DELETE CASCADE
);

CREATE TABLE rarity (
    gacha_system_id INTEGER NOT NULL,
    id SERIAL NOT NULL,
    name VARCHAR(50) NOT NULL,
    chance NUMERIC(7,2) NOT NULL CHECK (chance >= 0 AND chance <= 100),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (gacha_system_id, id),
    FOREIGN KEY (gacha_system_id)
        REFERENCES gacha_system(id)
        ON DELETE CASCADE
);

CREATE TABLE character (
   gacha_system_id INTEGER NOT NULL,
   id SERIAL NOT NULL,
   rarity_id INTEGER NOT NULL,
   name VARCHAR(100) NOT NULL,
   image_url TEXT,
   created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
   PRIMARY KEY (gacha_system_id, id),
   FOREIGN KEY (gacha_system_id)
       REFERENCES gacha_system(id)
       ON DELETE CASCADE,
   FOREIGN KEY (gacha_system_id, rarity_id)
       REFERENCES rarity(gacha_system_id, id)
       ON DELETE CASCADE
);