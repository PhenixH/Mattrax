-- DO NOT RUN THIS FILE. It is used along with sqlc to generate type safe Go from SQL

-------- User & Tenant

-- name: NewUser :exec
INSERT INTO users(upn, fullname, password, tenant_id) VALUES ($1, $2, $3, $4);

-- name: GetUser :one
SELECT upn, fullname, azuread_oid FROM users WHERE upn = $1 LIMIT 1;

-- name: GetUsersInTenant :many
SELECT upn, fullname, azuread_oid FROM users WHERE tenant_id = $1 LIMIT $2 OFFSET $3;
-- Once https://github.com/kyleconroy/sqlc/issues/778 is fixed change query to (including the ByQuery one): SELECT upn, fullname, azuread_oid FROM users INNER JOIN tenant_users ON users.upn = tenant_users.user_upn WHERE tenant_users.tenant_id = $1 UNION ALL SELECT upn, fullname, azuread_oid FROM users WHERE tenant_id = $1;

-- name: GetUsersInTenantByQuery :many
SELECT upn, fullname, azuread_oid FROM users WHERE tenant_id = $1 AND (upn || fullname || azuread_oid) LIKE $4 LIMIT $2 OFFSET $3;

-- name: GetUserSecure :one
SELECT fullname, password, mfa_token, tenant_id FROM users WHERE upn = $1 LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users WHERE upn=$1;

-- name: NewUserFromAzureAD :exec
INSERT INTO users(upn, fullname, azuread_oid) VALUES ($1, $2, $3);

-- name: NewTenant :one
INSERT INTO tenants(display_name, primary_domain) VALUES ($1, $2) RETURNING id;

-- name: ScopeUserToTenant :exec
INSERT INTO tenant_users(user_upn, tenant_id, permission_level) VALUES ($1, $2, $3);

-- name: GetUserPermissionLevelForTenant :one
SELECT permission_level FROM tenant_users WHERE user_upn = $1 AND tenant_id = $2;

-- name: GetUserTenants :many
SELECT id, display_name, primary_domain, description FROM tenants INNER JOIN tenant_users ON tenants.id = tenant_users.tenant_id WHERE tenant_users.user_upn = $1;

-- name: RemoveUserFromTenant :exec
DELETE FROM tenant_users WHERE user_upn=$1 AND tenant_id=$2;

-------- Group Actions

-- name: NewGroup :one
INSERT INTO groups(name, tenant_id) VALUES ($1, $2) RETURNING id;

-- name: GetGroups :many
SELECT id, name, description FROM groups WHERE tenant_id = $1 LIMIT $2 OFFSET $3;

-- name: GetGroup :one
SELECT id, name, description FROM groups WHERE id = $1 AND tenant_id = $2 LIMIT 1;

-------- Policy Actions

-- name: NewPolicy :one
INSERT INTO policies(name, tenant_id) VALUES ($1, $2) RETURNING id;

-- name: GetPolicies :many
SELECT id, name, description FROM policies WHERE tenant_id = $1 LIMIT $2 OFFSET $3;

-- name: GetPolicy :one
SELECT id, name, description FROM policies WHERE id = $1 AND tenant_id = $2 LIMIT 1;

-------- Device Actions

-- name: GetDevices :many
SELECT id, name, model FROM devices WHERE id = $1 AND tenant_id = $2 LIMIT $3 OFFSET $4;

-- name: GetDevice :one
SELECT id, name, description, model FROM devices WHERE id = $1 AND tenant_id = $2 LIMIT 1;

-- name: GetDeviceGroups :many
SELECT groups.id, groups.name FROM groups INNER JOIN group_devices ON group_devices.group_id=groups.id WHERE group_devices.device_id = $1;

-- name: GetDevicePolicies :many
SELECT id, name, description, policy_id, group_devices.group_id FROM policies INNER JOIN group_policies ON group_policies.policy_id = policies.id INNER JOIN group_devices ON group_devices.group_id=group_policies.group_id WHERE group_devices.device_id = $1;

-------- Certificates

-- name: GetRawCert :one
SELECT cert, key FROM certificates WHERE id = $1 LIMIT 1;

-- name: CreateRawCert :exec
INSERT INTO certificates(id, cert, key) VALUES ($1, $2, $3);