
CREATE TABLE env_managers (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  min_replicas INT NOT NULL,
  enabled BOOLEAN NOT NULL,
  ui_enabled BOOLEAN DEFAULT FALSE,
  last_update BIGINT NOT NULL,
  UNIQUE(name)
);

-- psql -d env_manager -U postgres -W --command "SELECT * FROM env_managers;"