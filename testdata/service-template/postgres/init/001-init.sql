CREATE TABLE smoke (
  id integer PRIMARY KEY,
  marker text NOT NULL
);

INSERT INTO smoke (id, marker) VALUES (1, 'zombie-postgres-ok');
