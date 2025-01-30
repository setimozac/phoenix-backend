
CREATE TABLE env_managers (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  min_replicas INT NOT NULL,
  enabled BOOLEAN NOT NULL,
  ui_enabled BOOLEAN DEFAULT FALSE,
  last_update BIGINT NOT NULL,
  namespace VARCHAR(255),
  cr_name VARCHAR(255),
  UNIQUE(name)
);

-- psql -d env_manager -U postgres -W --command "SELECT * FROM env_managers;"



-- events JSONB
-- INSERT INTO env_managers(name, min_replicas, enabled, last_update, namespace, cr_name, events) VALUES(...., '["BAU", "team1-test", "team2-test1", "team2-test2"]'::jsonb);
-- SELECT FROM env_managers WHERE events @> '["team2-test1"]'::jsonb;
-- CREATE INDEX events_gin_idx ON env_managers USING GIN (events);