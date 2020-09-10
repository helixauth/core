BEGIN;

CREATE TABLE IF NOT EXISTS tenants (
  id                    TEXT PRIMARY KEY,
  name                  TEXT DEFAULT NULL,
  picture               TEXT DEFAULT NULL,
  website               TEXT DEFAULT NULL,
  email                 TEXT DEFAULT NULL,
  email_provider        TEXT DEFAULT NULL,
  aws_region_id         TEXT DEFAULT NULL,
  aws_access_key_id     TEXT DEFAULT NULL,
  aws_secret_access_key TEXT DEFAULT NULL
);

ALTER TABLE tenants ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_policy ON tenants
    USING (id = current_setting('app.tenant_id'));

CREATE TABLE IF NOT EXISTS users (
  id                 TEXT PRIMARY KEY,
  tenant_id          TEXT NOT NULL REFERENCES tenants(id),
  name               TEXT DEFAULT NULL,
	nickname           TEXT DEFAULT NULL,
	preferred_username TEXT DEFAULT NULL,
	given_name         TEXT DEFAULT NULL,
	middle_name        TEXT DEFAULT NULL,
	family_name        TEXT DEFAULT NULL,
	email              TEXT DEFAULT NULL,
	email_verified     BOOLEAN DEFAULT NULL,
	zone_info          TEXT DEFAULT NULL,
	locale             TEXT DEFAULT NULL,
	address            TEXT DEFAULT NULL,
	phone_number       TEXT DEFAULT NULL,
	picture            TEXT DEFAULT NULL,
	website            TEXT DEFAULT NULL,
	gender             TEXT DEFAULT NULL,
	birthdate          TEXT DEFAULT NULL,
	is_blocked         BOOLEAN DEFAULT false,
  metadata           JSON NOT NULL,
  created_at         TIMESTAMP NOT NULL,
  updated_at         TIMESTAMP NOT NULL,
  last_active_at     TIMESTAMP DEFAULT NULL,
  CONSTRAINT unique_email_per_tenant UNIQUE (email, tenant_id)
);

ALTER TABLE users ENABLE ROW LEVEL SECURITY;

CREATE POLICY users_tenant_isolation_policy ON users
    USING (tenant_id = current_setting('app.tenant_id'));

CREATE TABLE IF NOT EXISTS clients (
  id                 TEXT PRIMARY KEY,
  tenant_id          TEXT NOT NULL REFERENCES tenants(id),
  name               TEXT NOT NULL,
  secret             TEXT NOT NULL,
  logo               TEXT DEFAULT NULL,
  url                TEXT DEFAULT NULL,
  description        TEXT DEFAULT NULL,
  privacy_policy     TEXT DEFAULT NULL,
  is_third_part      BOOLEAN DEFAULT true,
  authorized_domains TEXT[]
);

ALTER TABLE clients ENABLE ROW LEVEL SECURITY;

CREATE POLICY clients_tenant_isolation_policy ON clients
    USING (tenant_id = current_setting('app.tenant_id'));

CREATE TABLE IF NOT EXISTS sessions (
  id            TEXT PRIMARY KEY,
  tenant_id     TEXT NOT NULL REFERENCES tenants(id),
  user_id       TEXT NOT NULL REFERENCES users(id),
  client_id     TEXT NOT NULL REFERENCES clients(id),
  response_type TEXT NOT NULL,
  scope         TEXT NOT NULL,
  state         TEXT NOT NULL,
  nonce         TEXT NOT NULL,
  redirect_uri  TEXT NOT NULL,
  code          TEXT NOT NULL,
  created_at    TIMESTAMP NOT NULL,
  claimed_at    TIMESTAMP DEFAULT NULL,
  refreshed_at  TIMESTAMP DEFAULT NULL
);

ALTER TABLE sessions ENABLE ROW LEVEL SECURITY;

CREATE POLICY sessions_tenant_isolation_policy ON sessions
    USING (tenant_id = current_setting('app.tenant_id'));

CREATE TABLE IF NOT EXISTS passwords (
  id         TEXT PRIMARY KEY,
  tenant_id  TEXT NOT NULL REFERENCES tenants(id),
  user_id    TEXT NOT NULL REFERENCES users(id),
  hash       TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  retired_at TIMESTAMP DEFAULT NULL
);

ALTER TABLE passwords ENABLE ROW LEVEL SECURITY;

CREATE POLICY passwords_tenant_isolation_policy ON passwords
    USING (tenant_id = current_setting('app.tenant_id'));

-- WARNING: Update this user's password after bootstrapping the server
CREATE USER helix WITH ENCRYPTED PASSWORD 'password';

GRANT SELECT, INSERT, UPDATE, DELETE
  ON ALL TABLES IN SCHEMA public
  TO helix;

COMMIT;
