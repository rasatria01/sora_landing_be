package constants

const (
	AuthPasswordInvalidOrEmailNotFound = "Password wrong or Email not found"
	AuthLoginSuccess                   = "Login successful"
	AuthLogoutSuccess                  = "Logout successful"
	AuthRefreshTokenSuccess            = "Refresh token successful"
	AuthInvalidToken                   = "Invalid token"
)

type LoginType string

const (
	LoginTypeEmail LoginType = "email"
	LoginTypePhone LoginType = "phone"
)
