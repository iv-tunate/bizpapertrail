package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/iv-tunate/bizpapertrail/database"
	"github.com/iv-tunate/bizpapertrail/internal/auth"
	"github.com/iv-tunate/bizpapertrail/internal/cache"
	"github.com/iv-tunate/bizpapertrail/internal/models"
	"github.com/iv-tunate/bizpapertrail/internal/utils"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type RegisterParam struct {
    Name         *string `json:"name"          validate:"required"`
    Email        *string `json:"email"          validate:"required,email_regex"`
    Password     *string `json:"password"       validate:"required,min=8,password"`
    BusinessName *string `json:"business_name"  validate:"required"`
    PhoneNumber  *string `json:"phone_number"   validate:"required,phone_regex"`
}
func (h *Handler) RegisterAdmin(c echo.Context) error{
	var params RegisterParam
	ctx := c.Request().Context()

	if err := c.Bind(&params); err !=nil {
		h.Logger.ErrorContext(ctx, "bind failed", "error", err, "Handler:", "Adminr")
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid json body param", http.StatusBadRequest)
	}

	if validationErrors := utils.ValidateUserParams(params); validationErrors != nil {
		h.Logger.InfoContext(ctx, "Validation of registation parameter failed", "Details", map[string]any{
			"Handler": "RegisterAdmin",
			"Error(s)": validationErrors,
			"Email": *params.Email,
		})
    	return utils.ErrorResponse(c, http.StatusBadRequest, "validation failed", validationErrors)
	}

	hashedPassWord, err := bcrypt.GenerateFromPassword([]byte(*params.Password), bcrypt.DefaultCost)
	if err != nil{
		h.Logger.ErrorContext(ctx, "[ERROR]: An error occured while generating password hash", 
		 "Details", map[string]any{
			"Handler": "RegisterAdmin",
			"Error(s)": err,
			"Email": *params.Email,
			"status_code": 500,
		})
		return  utils.ErrorResponse(c, http.StatusInternalServerError, "An internal Server Error occued...Please try again", "Internal Server Error")
	}

	tx, err := h.Pool.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil{
		h.Logger.ErrorContext(ctx, "[ERROR] An error occured",
			"Error Details", map[string]any {
			"Handler": "RegisterAdmin", 
			"Error": err,
			"Email_Attempting Register Operation": *params.Email,
			},
		)
		return  utils.ErrorResponse(c, http.StatusInternalServerError, "An internal Server Error occued...Please try again", "Internal Server Error")
	}

	qtx := h.DB.WithTx(tx)

	password := string(hashedPassWord)

	userRow, err := qtx.CreateUser(ctx, database.CreateUserParams{
		Name:   strings.ToUpper(*params.Name),
		Email:  strings.ToLower(*params.Email),
		PhoneNumber: *params.PhoneNumber,
		BusinessName:   strings.ToUpper(*params.BusinessName),
		IsAdmin: true,
		Password: password,
	})

	if err != nil{
		statusCode, msg := utils.ParseDbError(err)
		h.Logger.ErrorContext(ctx, "[ERROR]: An error occured:", "db_error", msg, "email_attempted", *params.Email)
		return  utils.ErrorResponse(c, statusCode, msg, http.StatusText(statusCode))
	}

	token, err := utils.GenerateRandomToken(6)
	if err != nil{
		h.Logger.ErrorContext(ctx, "[ERROR] An error occured",
			"Error Details", map[string]any {
			"Handler": "RegisterUser", 
			"Error": err,
			"Email": userRow.Email,
			},
		)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "An internal Server Error occued...Please try again", "Internal Server Error")
	}
	
	expirationTime := time.Now().Add(5 * time.Minute)

	cache.StringCache.Set(utils.EmailVerificationKey(*params.Email), token, 5 * time.Minute)
	h.Logger.InfoContext(ctx, "Verification Token successsfully generated and cached", 
	"Details", map[string]any{
		"Email": *params.Email,
		"Token":token,
		"Expiration": expirationTime,
	})

	if err := tx.Commit(ctx); err != nil{
		h.Logger.ErrorContext(ctx, "[ERROR] An error occured while commiting the transaction",
		"Error Details", map[string]any {
		"Handler": "RegisterUser", 
		"Error": err,
		"context": ctx,
		"Email_Attempting register operation": *params.Email,
		})

		return utils.ErrorResponse(c, http.StatusInternalServerError, "An internal Server Error occued...Please try again", "Internal Server Error")
	}
	h.Logger.InfoContext(ctx,"[SUCCESS] RegisterUser: User profile created...", "user_data", userRow )
	return  utils.SuccessResponse(c, 200, "User Profile Created Successfully", userRow, nil)
}

func (h *Handler) RegisterUser(c echo.Context) error{
	var params RegisterParam
	ctx := c.Request().Context()

	if err := c.Bind(&params); err !=nil {
		h.Logger.ErrorContext(ctx, "bind failed", "error", err, "Handler:", "RegisterUser")
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid json body param", http.StatusBadRequest)
	}
	
	if validationErrors := utils.ValidateUserParams(params); validationErrors != nil {
		h.Logger.InfoContext(ctx, "Validation of registation parameter failed","Details", map[string]any{
			"Handler": "RegisterAdmin",
			"Error(s)": validationErrors,
			"Email": *params.Email,
		})
    	return utils.ErrorResponse(c, http.StatusBadRequest, "validation failed", validationErrors)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*params.Password), bcrypt.DefaultCost)
	if err != nil{
		h.Logger.ErrorContext(ctx, "[ERROR]: An error occured while generating password hash", 
		"error", err, "status_code", 500, "Details", map[string]any{
			"Handler": "RegisterUser",
			"Error(s)": err,
			"Email": *params.Email,
		})
		return  utils.ErrorResponse(c, http.StatusInternalServerError, "An internal Server Error occued...Please try again", "Internal Server Error")
	}

	tx, err := h.Pool.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil{
		h.Logger.ErrorContext(ctx, "[ERROR] An error occured",
			"Error Details", map[string]any {
			"Handler": "RegisterUser", 
			"Error": err,
			"Email_Attempting Register Operation": *params.Email,
			},
		)
		return  utils.ErrorResponse(c, http.StatusInternalServerError, "An internal Server Error occued...Please try again", "Internal Server Error")
	}

	qtx := h.DB.WithTx(tx)

	password := string(hashedPassword)

	userRow, err := qtx.CreateUser(ctx, database.CreateUserParams{
		Name:   strings.ToUpper(*params.Name),
		Email:   strings.ToLower(*params.Email),
		PhoneNumber: *params.PhoneNumber,
		BusinessName:   strings.ToUpper(*params.BusinessName),
		IsAdmin: false,
		Password: password,
	})

	if err != nil{
		statusCode, msg := utils.ParseDbError(err)
		h.Logger.ErrorContext(ctx, "[ERROR]: An error occured:", "db_error", msg, "email_attempted", *params.Email, "Handler", "Register User")
		return  utils.ErrorResponse(c, statusCode, msg, http.StatusText(statusCode) )
	}

	token, err := utils.GenerateRandomToken(6)
	if err != nil{
		h.Logger.ErrorContext(ctx, "[ERROR] An error occured",
		"Error Details", map[string]any {
		"Handler": "Register User", 
		"Error": err,
		"Email_Attempting register operation": *params.Email,
		},
		)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "An internal Server Error occued...Please try again", "Internal Server Error")
	}
	
	expirationTime := time.Now().Add(5 * time.Minute)

	cache.StringCache.Set(utils.EmailVerificationKey(*params.Email), token, 5 * time.Minute)
	h.Logger.InfoContext(ctx, "Verification Token successsfully generated and cached", 
	"Details", map[string]any{
		"Token":token,
		"Expiration": expirationTime,
	})

	if err := tx.Commit(ctx); err != nil{
		h.Logger.ErrorContext(ctx, "[ERROR] An error occured while commiting the transaction",
		"Error Details", map[string]any {
		"Handler": "Register User", 
		"Error": err,
		"context": ctx,
		"Email_Attempting register operation": *params.Email,
		})

		return utils.ErrorResponse(c, http.StatusInternalServerError, "An internal Server Error occued...Please try again", "Internal Server Error")
	}
	h.Logger.InfoContext(ctx,"[SUCCESS] RegisterUser: User profile created...", "user_data", userRow )
	return  utils.SuccessResponse(c, 200, "User Profile Created Successfully", userRow, nil)
}

func (h *Handler) VerifyUser(c echo.Context) error {

	type jsonBody struct{
		Email *string `json:"email" validate:"required,email"`
		Token *string `json:"token" validate:"required"`
		VerificationType *string `json:"verification_type" validate:"required"`
	}

	params := jsonBody{}

	ctx := c.Request().Context()
	if err := c.Bind(&params); err != nil{
		h.Logger.ErrorContext(ctx, "bind failed", "error", err, "Handler:", "RequestVerificationToken")
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid Json Body", http.StatusText(400));
	}
	cacheKey := utils.Tenary(*params.VerificationType == "Email", utils.EmailVerificationKey(*params.Email), utils.LoginVerificationKey(*params.Email))

	item := cache.StringCache.Get(cacheKey)
	if item == nil{
			h.Logger.ErrorContext(ctx, "Missing cache item", "Details",
			map[string]any {
			"Handler": "VerifyUser", 
			"UserEmail": *params.Email,
			},
		)
		return  utils.ErrorResponse(c, http.StatusBadRequest, "Invalid or expired verification token", "Bad Request")
	}

	tx, err := h.Pool.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil{
		h.Logger.ErrorContext(ctx, "[ERROR] An error occured",
			"Error Details", map[string]any {
			"Handler": "VerifyUser", 
			"Error": err,
			"UserEmail": *params.Email,
			},
		)
		return  utils.ErrorResponse(c, http.StatusInternalServerError, "An internal Server Error occued...Please try again", "Internal Server Error")
	}

	qtx := h.DB.WithTx(tx)

	if item.Value() != *params.Token{
		h.Logger.InfoContext(ctx, "Invalid Verification Token", 
		"Details:", map[string]any{
			"Handler": "Verify User",
			"User Email": *params.Email,
		})
		return  utils.ErrorResponse(c, 400, "Invalid verification token", http.StatusBadRequest)
	}

	var userRow database.VerifyUserEmailRow
	if *params.VerificationType == "Email"{

		userRow, err = qtx.VerifyUserEmail(ctx, strings.ToLower(*params.Email))
		if err != nil{
			statusCode, msg := utils.ParseDbError(err)
			h.Logger.ErrorContext(ctx, "[ERROR]: An error occured:", "db_error", msg, "email_attempted", *params.Email)
			return  utils.ErrorResponse(c, statusCode, msg,  http.StatusText(statusCode))
		}
	}

	claims := models.UserClaims{
		UserID: userRow.ID.String(),
		Email: userRow.Email,
		Role: utils.Tenary(userRow.IsAdmin, "Admin", "User"),
		Verified: userRow.Verified,
		Blacklisted: userRow.Blacklisted,
	}
	token, err := auth.GenerateJwtToken(claims)
	refreshToken, err := auth.GenerateRefreshToken(claims)
	if err != nil{
		h.Logger.ErrorContext(ctx, "[ERROR] An error occured while generating jwt token",
			"Error Details", map[string]any {
				"Handler": "VerifyUser", 
				"Error": err,
				"context": ctx,
				"UserEmail": *params.Email,
			})
		return utils.ErrorResponse(c, http.StatusInternalServerError, "An internal Server Error occued...Please try again", "Internal Server Error")
	}
	cache.StringCache.Set(utils.RefreshTokenKey(claims.UserID), refreshToken, 7 * 24 * time.Hour)

	if err := tx.Commit(ctx); err != nil {
		h.Logger.ErrorContext(ctx, "[ERROR] An error occured while commiting the transaction",
			"Error Details", map[string]any {
				"Handler": "VerifyUser", 
				"Error": err,
				"context": ctx,
				"UserEmail": *params.Email,
			})
		return utils.ErrorResponse(c, http.StatusInternalServerError, "An internal Server Error occued...Please try again", "Internal Server Error")
	}

	h.Logger.InfoContext(ctx, "Verification Successful", "claims:", claims)
	return utils.SuccessResponse(c, 200, "Verification Successful", userRow, map[string]any{
		"access_token" : token,
		"refresh_token" : refreshToken,
	})
}