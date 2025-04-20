CREATE TABLE IF NOT EXISTS roles (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  level int NOT NULL DEFAULT 0,
  description TEXT,
  created_at timestamptz NOT NULL DEFAULT (now()),
  updated_at timestamptz NOT NULL DEFAULT (now())
);

INSERT INTO
  roles (name, description, level)
VALUES
  (
    'user',
    'A user can create posts and comments',
    1
  );

INSERT INTO
  roles (name, description, level)
VALUES
  (
    'super',
    'An super can update and delete other users posts',
    3
  );


ALTER TABLE
  IF EXISTS users
ADD
  COLUMN role_id INT REFERENCES roles(id) DEFAULT 1;

UPDATE
  users
SET
  role_id = (
    SELECT
      id
    FROM
      roles
    WHERE
      name = 'user'
  );

ALTER TABLE
  users
ALTER COLUMN
  role_id DROP DEFAULT;

ALTER TABLE
  users
ALTER COLUMN
  role_id
SET
  NOT NULL;