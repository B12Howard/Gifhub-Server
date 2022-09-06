-------------------------------------------------------------------------------

\set ON_ERROR_STOP on


-------------------------------------------------------------------------------

-- Create the usage, user_files, users, & user_types tables which provide us a way to manage
-- user related data


-------------------------------------------------------------------------------

CREATE TABLE user_types (
    id SERIAL UNIQUE,
    file_size_limit SMALLINT,
    name VARCHAR(50)
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),,
    updated_at TIMESTAMP,

    PRIMARY KEY (id)
)

CREATE TABLE users (
    id SERIAL UNIQUE,
    uid INTEGER NOT NULL,
    user_type_id SMALLINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    udpatedat TIMESTAMP,
    deletedat TIMESTAMP,

    PRIMARY KEY (id),
    FOREIGN KEY (user_type_id)  REFERENCES user_types(id),
)

CREATE TABLE user_files (
    id SERIAL UNIQUE,
    uid INTEGER NOT NULL,
    url VARCHAR(2000) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()

    PRIMARY KEY (id),
    FOREIGN KEY (uid)  REFERENCES users(id),

)

CREATE TABLE usage (
    id SERIAL UNIQUE,
    uid INTEGER NOT NULL,
    duration INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()

    PRIMARY KEY (id),
    FOREIGN KEY (uid)  REFERENCES users(id),
)