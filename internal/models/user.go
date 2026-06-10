package models

type UserClaims struct {
	UserID string 
	Email string 
	Role string
	Verified bool
	Blacklisted bool 
}
