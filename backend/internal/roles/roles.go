package roles

import "fmt"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func (r Role) String() string {
	return string(r)
}

func (r Role) IsValid() bool {
	return r == RoleAdmin || r == RoleUser
}

func ParseRole(s string) (Role, error) {
	role := Role(s)
	if !role.IsValid() {
		return "", fmt.Errorf("invalid role: %s", s)
	}
	return role, nil
}

// Permission defines what actions are allowed
type Permission string

const (
	// Monitor permissions
	PermMonitorView    Permission = "monitor:view"
	PermMonitorDetails Permission = "monitor:details"

	// Review permissions
	PermReviewView   Permission = "review:view"
	PermReviewCreate Permission = "review:create"
	PermReviewEdit   Permission = "review:edit"
	PermReviewDelete Permission = "review:delete"

	// User management permissions
	PermUserView       Permission = "user:view"
	PermUserCreate     Permission = "user:create"
	PermUserEdit       Permission = "user:edit"
	PermUserDelete     Permission = "user:delete"
	PermUserChangeRole Permission = "user:changeRole"
)

// RolePermissions maps roles to their permissions
var RolePermissions = map[Role][]Permission{
	RoleUser: {
		PermMonitorView,
		PermMonitorDetails,
		PermReviewView,
		PermReviewCreate,
		PermUserView,
	},
	RoleAdmin: {
		PermMonitorView,
		PermMonitorDetails,
		PermReviewView,
		PermReviewCreate,
		PermReviewEdit,
		PermReviewDelete,
		PermUserView,
		PermUserCreate,
		PermUserEdit,
		PermUserDelete,
		PermUserChangeRole,
	},
}

// HasPermission checks if a role has a specific permission
func (r Role) HasPermission(p Permission) bool {
	permissions, ok := RolePermissions[r]
	if !ok {
		return false
	}
	for _, perm := range permissions {
		if perm == p {
			return true
		}
	}
	return false
}

// AllPermissions returns all permissions for a role
func (r Role) AllPermissions() []Permission {
	if perms, ok := RolePermissions[r]; ok {
		return perms
	}
	return []Permission{}
}
