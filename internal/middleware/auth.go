package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	parser "github.com/golang-jwt/jwt/v4"
	"github.com/iv-tunate/bizpapertrail/internal/auth"
	"github.com/iv-tunate/bizpapertrail/internal/cache"
	"github.com/iv-tunate/bizpapertrail/internal/utils"
	"github.com/labstack/echo"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc{
	return  func(c echo.Context) error{
		ctx := c.Request().Context()

		token, err := utils.ExtractJwtFromHeader(c)
		if err != nil{
			slog.InfoContext(ctx, "Empty or missing header or header value", "Details", map[string]any{
			"Location": "JWTMiddleware",
			})
			return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or Expired Token", http.StatusText(http.StatusUnauthorized))
	
		}

		claims := &auth.Claims{}
		jwtToken, _, err := parser.NewParser().ParseUnverified(token, claims)
		if err != nil{
			slog.ErrorContext(ctx, "Invalid Jwt token used for logout", "Error Details", map[string]any{
				"Location": "JWTMiddleware",
			})
			return  utils.ErrorResponse(c, http.StatusBadRequest, "Invalid authorization token", http.StatusBadRequest)
		}

		cacheKey := fmt.Sprintf("Blacklist:%s", jwtToken.Signature)
		
		if cache.BlacklistedTokensCache.Has(cacheKey){
			slog.ErrorContext(ctx, "An attempt to access a resource with revoked jwt token ", "Error Details", map[string]any{
				"Location": "JWTMiddleware",
			})
			return  utils.ErrorResponse(c, http.StatusForbidden, "Revoked Authorization Token", http.StatusText(http.StatusForbidden))
		}
		claims, err = auth.ValidateToken(token)
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