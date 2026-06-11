package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/iv-tunate/bizpapertrail/internal/auth"
	"github.com/iv-tunate/bizpapertrail/internal/cache"
	"github.com/iv-tunate/bizpapertrail/internal/models"
	"github.com/iv-tunate/bizpapertrail/internal/utils"
	"github.com/jellydator/ttlcache/v3"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
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

func(h *Handler) RequestVerificationToken(c echo.Context) error{
	type jsonBody struct{
		Email *string `json:"email" validate:"required,email"`
		VerificationType *string `json:"verification_type" validate:"required"`
	}

	params := &jsonBody{}
	ctx := c.Request().Context()

	if err := c.Bind(params); err != nil{
		h.Logger.ErrorContext(ctx, "bind failed", "error", err, "Handler:", "RequestVerificationToken")
		return utils.ErrorResponse(c, http.StatusBadRequest, "An Error occured", http.StatusBadRequest);
	}

	userExists, err := h.DB.CheckUserExistsViaEmail(ctx, *params.Email)
	if err != nil{
		statusCode, msg := utils.ParseDbError(err)
		h.Logger.ErrorContext(ctx, "An error occured:", "db_error", msg, "email_attempted", *params.Email, "Handler", "RequestVerificationToken")
		return  utils.ErrorResponse(c, statusCode, msg, http.StatusText(statusCode))
	}

	if !userExists{
		h.Logger.InfoContext(ctx, "Attempt to request a verification token by nonexistent email", "Details", map[string]any{
			"Handler": "RequestVerificationToken",
			"Email": *params.Email,
		})

		return utils.ErrorResponse(c, 400, "user email does not exist", "Bad Request")
	}

	token, err := utils.GenerateRandomToken(6)
	if err != nil{
		h.Logger.ErrorContext(ctx, "An error occured",
			"Error Details", map[string]any {
			"Handler": "RequestVerificationToken", 
			"Error": err,
			"Email": *params.Email,
			},
		)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "An internal server error occured...Please try again", "Internal Server Error")
	}
	
	expirationTime := time.Now().Add(5 * time.Minute)
	cacheKey := utils.Tenary(*params.VerificationType == "Email", utils.EmailVerificationKey(*params.Email), utils.LoginVerificationKey(*params.Email))
	
	cache.StringCache.Set(cacheKey, token, ttlcache.DefaultTTL)

	h.Logger.InfoContext(ctx, "Verification Token successsfully generated and cached", 
	"Details", map[string]any{
		"Email": *params.Email,
		"Token":token,
		"Expiration": expirationTime,
	})
	return  utils.SuccessResponse(c, 200, "Verification Code has been sent", nil, nil)
}

func (h *Handler) Login(c echo.Context) error {
	type body struct{
		Email *string `json:"email" validate:"required,email"`
		Password *string `json:"password" validate:"required"`
	}

	param := &body{}
	ctx := c.Request().Context()

	if err := c.Bind(param); err != nil{
		h.Logger.ErrorContext(ctx, "bind failed", "error", err, "Handler:", "Login")
		return utils.ErrorResponse(c, http.StatusBadRequest, "An Error occured", http.StatusBadRequest);
	}

	userRow, err := h.DB.GetUserDetails(ctx, *param.Email)
	if err != nil{
		_, msg := utils.ParseDbError(err)
		h.Logger.ErrorContext(ctx, "An error occured:", "db_error", msg, "email_attempted", *param.Email, "Handler", "Login")
		return  utils.ErrorResponse(c, 400, "Incorrect email or password", http.StatusText(400))
	}

	err = bcrypt.CompareHashAndPassword([]byte(userRow.Password), []byte(*param.Password))
	if err != nil{
		h.Logger.ErrorContext(ctx, "Wrong password", "email_attempted", *param.Email, "Handler", "Login", "err:", err)
		return  utils.ErrorResponse(c, 400, "Incorrect email or password", http.StatusText(400))
	}
	
	h.Logger.InfoContext(ctx, "Operation Successful", "Details", map[string]any{
		"Handler": "Login",
		"Message": "User credentials were valid for login",
		"Email": userRow.Email,
	})

	return  utils.SuccessResponse(c, 200, "Operation Successful", nil, nil)
}

func (h *Handler) Logout(c echo.Context) error {
	user_id := c.Get(utils.UserIdKey).(string)
	
	token, err := utils.ExtractJwtFromHeader(c)
	ctx := c.Request().Context()
	if err != nil{
		h.Logger.InfoContext(ctx, "Empty or missing header or header value", "Details", map[string]any{
		"Location": "JWTMiddleware",
		})
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or Expired Token", http.StatusText(http.StatusUnauthorized))
	}

	claims := &auth.Claims{}

	jwtToken, _, err := jwt.NewParser().ParseUnverified(token, claims)
	if err != nil{
		h.Logger.ErrorContext(ctx, "Invalid Jwt token used for logout", "Error Details", map[string]any{
			"Location": "Logout",
			"UserId": user_id,
		})
		return  utils.ErrorResponse(c, http.StatusBadRequest, "Invalid authorization token", http.StatusBadRequest)
	}

	timeLeft := time.Until(claims.ExpiresAt.Time)

	if timeLeft > 0{
		cacheKey := fmt.Sprintf("Blacklist:%s", jwtToken.Signature)
		cache.BlacklistedTokensCache.Set(cacheKey, true, timeLeft)
	}

	return utils.SuccessResponse(c, 201, "Logout Successful", nil, nil)
}