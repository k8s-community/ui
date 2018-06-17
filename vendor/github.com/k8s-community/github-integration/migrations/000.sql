CREATE TABLE installations (
  id              SERIAL PRIMARY KEY,
  username        VARCHAR(128) NOT NULL UNIQUE,
  installation_id INTEGER      NOT NULL UNIQUE,

  created_at      TIMESTAMP    NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE builds (
  id              SERIAL PRIMARY KEY,
  uuid            VARCHAR(512) NOT NULL UNIQUE,
  username        VARCHAR(128) NOT NULL,
  repository      VARCHAR(128) NOT NULL,
  commit          VARCHAR(512) NOT NULL,
  passed          BOOL NOT NULL DEFAULT FALSE,
  log             TEXT NOT NULL DEFAULT '',

  created_at      TIMESTAMP    NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMP    NOT NULL DEFAULT NOW()
)
