
CREATE TABLE env_managers (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  min_replicas INT NOT NULL,
  enable BOOLEAN NOT NULL
);