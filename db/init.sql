CREATE TABLE users (
  id          SERIAL PRIMARY KEY,
  source      VARCHAR(128) NOT NULL,
  name        VARCHAR(128) NOT NULL,
  session_id  VARCHAR(256) DEFAULT NULL UNIQUE,
  session_data TEXT DEFAULT NULL,

  token TEXT DEFAULT NULL,
  ca_crt TEXT DEFAULT NULL,

  created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMP NOT NULL DEFAULT NOW(),

  CONSTRAINT u_source_name UNIQUE (source, name)
);
