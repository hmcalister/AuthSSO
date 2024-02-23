CREATE TABLE users (
    uuid text PRIMARY KEY,
    username text NOT NULL,
    FOREIGN KEY (uuid) REFERENCES authenticationData(uuid)
)

CREATE TABLE authenticationData (
    uuid text PRIMARY KEY,
    hashedPassword text NOT NULL,
    salt text NOT NULL
)