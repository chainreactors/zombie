USE zombie;

CREATE TABLE smoke (
  id INT PRIMARY KEY,
  marker VARCHAR(64) NOT NULL
);

INSERT INTO smoke (id, marker) VALUES (1, 'zombie-mysql-ok');

CREATE TABLE app_credentials (
  id INT PRIMARY KEY,
  username VARCHAR(64) NOT NULL,
  password_hash VARCHAR(128) NOT NULL,
  api_token VARCHAR(128) NOT NULL
);

INSERT INTO app_credentials (id, username, password_hash, api_token)
VALUES (1, 'demo', 'hash', 'token');
