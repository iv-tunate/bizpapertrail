package handlers

import (
	"net/http"
	"time"

	"github.com/iv-tunate/bizpapertrail/internal/auth"
	"github.com/iv-tunate/bizpapertrail/internal/cache"
	"github.com/iv-tunate/bizpapertrail/internal/models"
	"github.com/iv-tunate/bizpapertrail/internal/utils"
	"github.com/labstack/echo"
)


func (h *Handler) RefreshJwtToken(c echo.Context) error{
 	type body struct{
		RefreshToken *string `json:"refresh_token"`
	}

	param := &body{}
	ctx := c.Request().Context()

	user_id := c.Get(utils.UserIdKey).(string)
	if user_id == ""{
		h.Logger.InfoContext(ctx, "Missing user id", "Details", map[string]any{
			"handler":"RefreshJwtToken",
		})
		return utils.ErrorResponse(c, 400, "An error occured", "Bad Request")
	}

	if err := c.Bind(param); err != nil {
        return utils.ErrorResponse(c, 400, "invalid request body", err)
    }

	item:= cache.StringCache.Get(utils.RefreshTokenKey(user_id))
		if item == nil{
			h.Logger.ErrorContext(ctx, "Missing cache item", "Details",
			map[string]any {
			"Handler": "RefreshJwtToken", 
			"StatusCode": 400,
			"UserId": user_id,
			},
		)
		return  utils.ErrorResponse(c, http.StatusBadRequest, "Invalid or expired verification token", "Bad Request")
	}

	claims, err := auth.ValidateRefreshToken(item.Value())
	if err != nil {
        return utils.ErrorResponse(c, 401, "invalid refresh token", err)
    }
	
	cache.StringCache.Delete(utils.RefreshTokenKey(user_id))

	userClaims := models.UserClaims{
		UserID: claims.UserID,
		Email: claims.Email,
		Role: claims.Role,
		Verified: claims.Verified,
		Blacklisted: claims.Blacklisted,
	}
	accessToken, err := auth.GenerateJwtToken(userClaims)

	if err != nil {
		h.Logger.InfoContext(ctx, "An error occured while generating JWT token", "Details", map[string]any{
			"Handler": "RefreshJwtToken",
			"claims": claims,
		})
        return utils.ErrorResponse(c, 500, "failed to generate token", err)
    }

	refreshToken, err := auth.GenerateRefreshToken(userClaims)
		if err != nil {
		h.Logger.InfoContext(ctx, "An error occured while generating JWT token", "Details", map[string]any{
			"Handler": "RefreshJwtToken",
			"claims": claims,
		})
		cache.StringCache.Set(utils.RefreshTokenKey(claims.UserID), refreshToken, 7 * 24 * time.Hour)
    }

	h.Logger.InfoContext(ctx, "Jwt Token generated successfully", "Details", map[string]any{
		"Handler": "RefreshJwtToken",
		"claims": claims,
	})
	return utils.SuccessResponse(c, 200, "Operation Successful", map[string]any{
		"access_token" : accessToken,
		"refresh_token" : refreshToken,
	}, nil)
}