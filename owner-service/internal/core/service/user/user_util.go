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
