
CREATE TABLE env_managers (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  min_replicas INT NOT NULL,
  enabled BOOLEAN NOT NULL,
  ui_enabled BOOLEAN DEFAULT FALSE,
  last_update BIGINT NOT NULL
);

-- INSERT INTO env_managers VALUES(1, 'service1', 3, true, false, 1737328284)
-- INSERT INTO env_managers VALUES(2, 'service2', 1, true, true, 1737328288)