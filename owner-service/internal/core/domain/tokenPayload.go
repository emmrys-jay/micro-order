package domain

// Claims is an entity that represents the payload of the token
type Claims struct {
	Email  string
	Issuer string
	ID     string
	// jwt.RegisteredClaims
}
