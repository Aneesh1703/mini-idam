package rbac

import "errors"

// Define roles
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
	RoleSuper = "superadmin"
)

// Map roles to allowed permissions
var rolePermissions = map[string][]string{
	RoleUser:  {"read_own"},
	RoleAdmin: {"read_own", "read_all", "manage_users"},
	RoleSuper: {"read_own", "read_all", "manage_users", "manage_vault"},
}

// CheckPermission returns true if the given role has the permission
func CheckPermission(role, permission string) bool {
	perms, ok := rolePermissions[role]
	if !ok {
		return false
	}

	for _, p := range perms {
		if p == permission {
			return true
		}
	}
	return false
}

// Middleware-like helper: returns error if role lacks permission
func EnforcePermission(role, permission string) error {
	if !CheckPermission(role, permission) {
		return errors.New("permission denied")
	}
	return nil
}
