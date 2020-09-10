BEGIN;

DROP POLICY IF EXISTS passwords_tenant_isolation_policy ON passwords;

DROP TABLE IF EXISTS passwords;

DROP POLICY IF EXISTS sessions_tenant_isolation_policy ON passwords;

DROP TABLE IF EXISTS sessions;

DROP POLICY IF EXISTS clients_tenant_isolation_policy ON passwords;

DROP TABLE IF EXISTS clients;

DROP POLICY IF EXISTS users_tenant_isolation_policy ON passwords;

DROP TABLE IF EXISTS users;

DROP POLICY IF EXISTS tenant_isolation_policy ON tenants;

DROP TABLE IF EXISTS tenants;

DROP OWNED BY helix;

DROP USER IF EXISTS helix;

COMMIT;
