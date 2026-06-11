package utils

import (
	"errors"
	"strings"

	"github.com/labstack/echo"
)

func ExtractJwtFromHeader(c echo.Context)(string, error){
	headerVal := c.Request().Header.Get("Authorization")

	if headerVal == "" || !strings.HasPrefix(headerVal, "Bearer"){
		return  "", errors.New("Empty or missing header or header value")
	}

	token := strings.TrimPrefix(headerVal, "Bearer")
	return  token, nil
}