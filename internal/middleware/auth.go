package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/iv-tunate/bizpapertrail/internal/auth"
	"github.com/iv-tunate/bizpapertrail/internal/utils"
	"github.com/labstack/echo"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc{
	return  func(c echo.Context) error{
		authHeader := c.Request().Header.Get("Authorization")
		ctx := c.Request().Context()
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer"){
			slog.InfoContext(ctx, "Missing or invalid authorization header")
			return utils.ErrorResponse(c, http.StatusUnauthorized, "Missing or invalid authorization header", http.StatusText(http.StatusUnauthorized))
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer")
		claims, err := auth.ValidateToken(tokenStr)
		if err != nil{
			slog.InfoContext(ctx, "Invalid or Expired Token", "Details", map[string]any{
				"claims": claims,
			})
			return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or Expired Token", http.StatusText(http.StatusUnauthorized))
		}

		if !claims.Verified{
			slog.InfoContext(ctx, "Attempt to access resource by unverified user", "claims:", claims)
			return utils.ErrorResponse(c, http.StatusForbidden, "User email unverified", http.StatusText(http.StatusForbidden))
		}

		if !claims.Blacklisted{
			slog.InfoContext(ctx, "Attempt to access resource by blacklisted user", "claims:", claims)
			return utils.ErrorResponse(c, http.StatusForbidden, "Blacklisted user", http.StatusText(http.StatusForbidden))
		}
		c.Set(utils.UserIdKey, claims.UserID)
		c.Set(utils.UserRoleKey, claims.Role)

		return  next(c)
	}
}

func AdminMiddleWare(next echo.HandlerFunc) echo.HandlerFunc{
	return  func(c echo.Context) error{
		role, ok := c.Get(utils.UserRoleKey).(string)
		if !ok || role != "Admin"{
				slog.InfoContext(c.Request().Context(), "Forbidden ttempt to access admin resource by non admin user")
			return utils.ErrorResponse(c, http.StatusForbidden, "Forbidden Access", http.StatusText(http.StatusForbidden))
		}
		return  next(c)
	}
}