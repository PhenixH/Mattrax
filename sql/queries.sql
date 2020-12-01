-- DO NOT RUN THIS FILE. It is used along with sqlc to generate type safe Go from SQL

-------- Tenant

-- name: GetTenant :one
SELECT display_name, primary_domain, email, phone FROM tenants WHERE id = $1 LIMIT 1;

-- name: GetTenantDomains :many
SELECT domain, linking_code, verified FROM tenant_domains WHERE tenant_id = $1;

-- name: GetTenantDomain :one
SELECT linking_code, verified FROM tenant_domains WHERE domain=$1 AND tenant_id=$2 LIMIT 1;

-- name: AddDomainToTenant :one
INSERT INTO tenant_domains(tenant_id, domain) VALUES ($1, $2) RETURNING linking_code;

-- name: UpdateDomain :exec
UPDATE tenant_domains SET verified=$3 WHERE domain=$1 AND tenant_id=$2;

-- name: DeleteDomain :exec
DELETE FROM tenant_domains WHERE domain=$1 AND tenant_id=$2;

-------- User

-- name: NewUser :exec
INSERT INTO users(upn, fullname, password, tenant_id) VALUES ($1, $2, $3, $4);

-- name: GetUser :one
SELECT upn, fullname, disabled, azuread_oid FROM users WHERE upn = $1 LIMIT 1;

-- name: GetUsersInTenant :many
SELECT upn, fullname, azuread_oid FROM users WHERE tenant_id = $1 LIMIT $2 OFFSET $3;
-- Once https://github.com/kyleconroy/sqlc/issues/778 is fixed change query to (including the ByQuery one): SELECT upn, fullname, azuread_oid FROM users INNER JOIN tenant_users ON users.upn = tenant_users.user_upn WHERE tenant_users.tenant_id = $1 UNION ALL SELECT upn, fullname, azuread_oid FROM users WHERE tenant_id = $1;

-- name: GetUsersInTenantByQuery :many
SELECT upn, fullname, azuread_oid FROM users WHERE tenant_id = $1 AND (upn || fullname || azuread_oid) LIKE $4 LIMIT $2 OFFSET $3;

-- name: GetUserSecure :one
SELECT fullname, password, mfa_token, tenant_id FROM users WHERE upn = $1 LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users WHERE upn=$1;

-- name: DeleteUserInTenant :exec
DELETE FROM users WHERE upn=$1 AND tenant_id=$2;

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

-------- Object Actions

-- name: GetObject :one
SELECT filename, data FROM objects WHERE id = $1 AND tenant_id = $2 LIMIT 1;

-- name: CreateObject :exec
INSERT INTO objects(tenant_id, filename, data) VALUES ($1, $2, $3) RETURNING id;

-- name: UpdateObject :exec
UPDATE objects SET filename=$3, data=$4 WHERE id=$1 AND tenant_id=$2;

-------- Group Actions

-- name: NewGroup :one
INSERT INTO groups(name, tenant_id) VALUES ($1, $2) RETURNING id;

-- name: GetGroups :many
SELECT id, name, description FROM groups WHERE tenant_id = $1 LIMIT $2 OFFSET $3;

-- name: GetGroup :one
SELECT id, name, description FROM groups WHERE id = $1 AND tenant_id = $2 LIMIT 1;

-- name: DeleteGroup :exec
DELETE FROM groups WHERE id = $1 AND tenant_id = $2;

-- name: GetDevicesInGroup :many
SELECT device_id FROM group_devices WHERE group_id = $1 LIMIT $2 OFFSET $3;

-- name: AddDevicesToGroup :exec
INSERT INTO group_devices(group_id, device_id) VALUES ($1, $2);

-- name: GetPoliciesInGroup :many
SELECT policy_id FROM group_policies WHERE group_id = $1 LIMIT $2 OFFSET $3;

-- name: AddPoliciesToGroup :exec
INSERT INTO group_policies(group_id, policy_id) VALUES ($1, $2);

-------- Policy Actions

-- name: NewPolicy :one
INSERT INTO policies(name, type, tenant_id) VALUES ($1, $2, $3) RETURNING id;

-- name: GetPolicies :many
SELECT id, name, type, description FROM policies WHERE tenant_id = $1 LIMIT $2 OFFSET $3;

-- name: GetPolicy :one
SELECT name, type, payload, description FROM policies WHERE id = $1 AND tenant_id = $2 LIMIT 1;

-- name: DeletePolicy :exec
DELETE FROM policies WHERE id = $1 AND tenant_id = $2;

-------- Application Actions

-- name: NewApplication :one
INSERT INTO applications(name, tenant_id) VALUES ($1, $2) RETURNING id;

-- name: GetApplications :many
SELECT id, name, publisher FROM applications WHERE tenant_id = $1 LIMIT $2 OFFSET $3;

-- name: GetApplication :one
SELECT name, description, publisher FROM applications WHERE id = $1 AND tenant_id = $2 LIMIT 1;

-- name: GetApplicationTargets :many
SELECT msi_file, store_id FROM application_target WHERE app_id = $1 AND tenant_id = $2;

-- name: UpdateApplication :exec
UPDATE applications SET name=COALESCE($3, name), description=COALESCE($4, description), publisher=COALESCE($5, publisher) WHERE id = $1 AND tenant_id=$2;

-- name: DeleteApplication :exec
DELETE FROM applications WHERE id = $1 AND tenant_id = $2;

-------- Device Actions

-- name: GetDevices :many
SELECT id, name, model FROM devices WHERE tenant_id = $1 LIMIT $2 OFFSET $3;

-- name: GetDevice :one
SELECT id, protocol, name, description, state, owner, azure_did, enrolled_at, model FROM devices WHERE id = $1 AND tenant_id = $2 LIMIT 1;

-- name: GetDeviceGroups :many
SELECT groups.id, groups.name FROM groups INNER JOIN group_devices ON group_devices.group_id=groups.id WHERE group_devices.device_id = $1;

-- name: GetDevicePolicies :many
SELECT id, name, description, policy_id, group_devices.group_id FROM policies INNER JOIN group_policies ON group_policies.policy_id = policies.id INNER JOIN group_devices ON group_devices.group_id=group_policies.group_id WHERE group_devices.device_id = $1;

-------- Certificates

-- name: GetRawCert :one
SELECT cert, key FROM certificates WHERE id = $1 LIMIT 1;

-- name: CreateRawCert :exec
INSERT INTO certificates(id, cert, key) VALUES ($1, $2, $3);