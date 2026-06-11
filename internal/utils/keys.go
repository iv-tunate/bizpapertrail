package utils

import "fmt"

var UserIdKey = "user_id"
var UserEmailKey = "email"
var UserRoleKey = "user_role"

func EmailVerificationKey(email string) string{
	return fmt.Sprintf("Email_Verification:%s", email)
}

func LoginVerificationKey(email string) string{
		return fmt.Sprintf("Login_Verification:%s", email)
}
func RefreshTokenKey(user_id string) string{
	return fmt.Sprintf("RefreshToken:%s", user_id)
}