CREATE TABLE users (
    uuid text PRIMARY KEY,
    username text NOT NULL UNIQUE,
    FOREIGN KEY (uuid) REFERENCES authenticationData(uuid)
)

CREATE TABLE authenticationData (
    uuid text PRIMARY KEY,
    hashed_password text NOT NULL,
    salt text NOT NULL
)