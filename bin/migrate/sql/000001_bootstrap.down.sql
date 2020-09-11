BEGIN;

DROP POLICY IF EXISTS email_verifications_tenant_isolation_policy ON email_verifications;

DROP TABLE IF EXISTS email_verifications;

DROP POLICY IF EXISTS sessions_tenant_isolation_policy ON sessions;

DROP TABLE IF EXISTS sessions;

DROP POLICY IF EXISTS users_tenant_isolation_policy ON users;

DROP TABLE IF EXISTS users;

DROP POLICY IF EXISTS clients_tenant_isolation_policy ON clients;

DROP TABLE IF EXISTS clients;

DROP POLICY IF EXISTS admins_tenant_isolation_policy ON admins;

DROP TABLE IF EXISTS admins;

DROP POLICY IF EXISTS tenant_isolation_policy ON tenants;

DROP TABLE IF EXISTS tenants;

DROP OWNED BY helix;

DROP USER IF EXISTS helix;

COMMIT;
