package constants

type (
	UserRole    string
	UserStatus  string
	UserRoleMap map[UserRole]bool
)

const (
	UserRoleAdmin UserRole = "Admin"
	UserRoleUser  UserRole = "User"
)

const (
	UserStatusActive = "Active"
)

func (receiver UserRole) IsValidEnum() bool {
	switch receiver {
	case UserRoleAdmin, UserRoleUser:
		return true
	default:
		return false
	}
}

func (receiver UserStatus) IsValidEnum() bool {
	switch receiver {
	case UserStatusActive:
		return true
	default:
		return false
	}
}
