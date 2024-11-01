BEGIN;

CREATE TABLE IF NOT EXISTS players (
    id UUID PRIMARY KEY,
    nickname VARCHAR(64) NOT NULL UNIQUE,
    email VARCHAR(128) NOT NULL UNIQUE,
    firstname VARCHAR(64) NOT NULL,
    lastname VARCHAR(64) NOT NULL,
    country VARCHAR(64) NOT NULL,
    usrpwd TEXT NOT NULL,
    date_created TIMESTAMP,
    date_updated TIMESTAMP
);

COMMIT;