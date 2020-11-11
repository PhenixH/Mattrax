CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Note: This function will have collisions after huge amounts of entities are created. For now this is fine but in future will need to be fixed.
CREATE OR REPLACE FUNCTION short_uuid() RETURNS TEXT AS $$
BEGIN
return replace(replace(encode(gen_random_bytes(6), 'base64'), '/', '_'), '+', '-');
end;
$$ language 'plpgsql';

CREATE TABLE tenants (
	id TEXT PRIMARY KEY DEFAULT short_uuid(),
    display_name TEXT NOT NULL,
    primary_domain TEXT UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE users (
    upn TEXT PRIMARY KEY,
    fullname TEXT NOT NULL,
    password TEXT,
    mfa_token TEXT,
    azuread_oid TEXT UNIQUE,
    tenant_id TEXT REFERENCES tenants(id)
);

CREATE TYPE user_permission_level AS ENUM ('user', 'administrator');

CREATE TABLE tenant_users (
	user_upn TEXT REFERENCES users(upn) NOT NULL,
    tenant_id TEXT REFERENCES tenants(id) NOT NULL,
    permission_level user_permission_level NOT NULL DEFAULT 'user',
    PRIMARY KEY (user_upn, tenant_id)
);

CREATE TYPE management_protocol AS ENUM ('windows', 'agent', 'apple');
CREATE TYPE device_state AS ENUM ('deploying', 'managed', 'user_unenrolled', 'missing');

CREATE TABLE devices (
    id TEXT PRIMARY KEY DEFAULT short_uuid(),
    tenant_id TEXT REFERENCES tenants(id) NOT NULL,
    protocol management_protocol NOT NULL,
    udid TEXT NOT NULL,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    model TEXT,
    state device_state NOT NULL,
    owner TEXT REFERENCES users(upn),
    azure_did TEXT UNIQUE,
    enrolled_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    UNIQUE(protocol, udid)
);

CREATE TABLE policies (
    id TEXT PRIMARY KEY DEFAULT short_uuid(),
    tenant_id TEXT REFERENCES tenants(id) NOT NULL,
    name TEXT UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE groups (
    id TEXT PRIMARY KEY DEFAULT short_uuid(),
    tenant_id TEXT REFERENCES tenants(id) NOT NULL,
    name TEXT UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE group_devices (
    group_id TEXT REFERENCES groups(id) NOT NULL,
    device_id TEXT REFERENCES devices(id) NOT NULL,
    PRIMARY KEY (group_id, device_id)
);

CREATE TABLE group_policies (
    group_id TEXT REFERENCES groups(id),
    policy_id TEXT REFERENCES policies(id),
    PRIMARY KEY (group_id, policy_id)
);

CREATE TABLE certificates (
    id TEXT PRIMARY KEY,
    cert BYTEA,
    key BYTEA
);