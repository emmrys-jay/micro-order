package user

import "fmt"

// UserRoleToString maps UserRole enum values to their string representations.
func UserRoleToString(userRole UserRole) string {
	switch userRole {
	case UserRole_USER_ROLE_ADMIN:
		return "admin"
	case UserRole_USER_ROLE_USER:
		return "user"
	default:
		return "UNKNOWN"
	}
}

// StringToUserRole maps string representations to UserRole enum values.
func StringToUserRole(userRoleStr string) (UserRole, error) {
	switch userRoleStr {
	case "admin":
		return UserRole_USER_ROLE_ADMIN, nil
	case "user":
		return UserRole_USER_ROLE_USER, nil
	default:
		return UserRole_USER_ROLE_UNSPECIFIED, fmt.Errorf("unknown user type: %s", userRoleStr)
	}
}

type UserUpdateForQueue struct {
	Id        string   `json:"id,omitempty"`
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
	Email     string   `json:"email,omitempty"`
	Password  string   `json:"password,omitempty"`
	Phone     string   `json:"phone,omitempty"`
	Role      UserRole `json:"role,omitempty"`
	IsActive  bool     `json:"is_active,omitempty"`
	CreatedAt string   `json:"created_at,omitempty"`
	UpdatedAt string   `json:"updated_at,omitempty"`
	DeletedAt string   `json:"deleted_at,omitempty"`
}
