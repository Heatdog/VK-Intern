CREATE TYPE users_role AS enum ('Administrator', 'User');

CREATE TABLE IF NOT EXISTS Users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role users_role NOT NULL DEFAULT 'User',
    CHECK (LENGTH(login) > 3 and LENGTH(password) > 3)
);

INSERT INTO Users VALUES ('Admin', 'Admin', 'Admin');
INSERT INTO Users VALUES ('VK Senior', '678', 'Admin');
INSERT INTO Users VALUES ('John', '123');
INSERT INTO Users VALUES ('David', '123456');
INSERT INTO Users VALUES ('Peter', '890');