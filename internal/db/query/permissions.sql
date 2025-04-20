SELECT p.name
FROM permissions p
JOIN roles_permissions rp ON p.id = rp.permission_id
JOIN roles r ON r.id = rp.role_id
JOIN users u ON u.role_id = r.id
WHERE u.id = $1;

-- name: GetPermissionsByRoleID :many
SELECT p.* FROM permissions p
JOIN roles_permissions rp ON rp.permission_id = p.id
WHERE rp.role_id = $1;
