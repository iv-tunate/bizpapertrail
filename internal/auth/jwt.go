package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/iv-tunate/bizpapertrail/internal/models"
)

type Claims struct{
	UserID string `json:"user_id"`
	Email string `json:"email"`
	Role string `json:"role"`
	Verified bool `json:"Verified"`
	Blacklisted bool `json:"Blacklisted"`
	jwt.RegisteredClaims

}

var environment = os.Getenv("ENV")
var issuer = os.Getenv(fmt.Sprintf("%s_ISSUER", environment))

func GenerateJwtToken(claims models.UserClaims)(string, error){
	userClaims := &Claims{
		UserID: claims.UserID,
		Email: claims.Email,
		Role: claims.Role,
		Verified: claims.Verified,
		Blacklisted: claims.Blacklisted,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaims)
	return  token.SignedString([]byte(os.Getenv("JUUT_SICKRIT")))
}

func ValidateToken(tokenStr string) (*Claims, error){
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func (t *jwt.Token) (any, error) {
		return []byte(os.Getenv("JUUT_SICKRIT")), nil
	})

	if err != nil || !token.Valid{
		return nil, err
	}
	return  claims, nil
}

func GenerateRefreshToken(claims models.UserClaims) (string, error){
	userClaims := &Claims{
		UserID: claims.UserID,
		Email: claims.Email,
		Role: claims.Role,
		Verified: claims.Verified,
		Blacklisted: claims.Blacklisted,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	return  token.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
}

func ValidateRefreshToken(tokenStr string) (*Claims, error){
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func (t *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})

	if err != nil || !token.Valid{
		return nil, err
	}
	return  claims, nil
}