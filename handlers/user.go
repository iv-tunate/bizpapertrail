package handlers

import (
	"net/http"
	"github.com/iv-tunate/bizpapertrail/database"
	"github.com/iv-tunate/bizpapertrail/utils"
	"github.com/labstack/echo"
)

type RegisterParam struct{
	Name *string `json:"name" validate:"required"`
	Email *string `json:"email" validate:"required, email_regex"`
	Password *string `json:"password" validate:"required, min=8, password_regex"`
	BusinessName *string `json:"business_name" validate:"requred"`
	PhoneNumber *string `json:"phone_number" validate:"required, phone_regex"`
	IsAdmin *bool `json:"is_admin"`
}

func (h *Handler) RegisterUser(c echo.Context)error{
	var params RegisterParam
	
	if err := c.Bind(&params); err !=nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid json body param", http.StatusBadRequest)
	}
	var paramContainer interface{} = params

	if registerParam, ok := paramContainer.(RegisterParam); ok{
		isValid := utils.ValidateUserParams(registerParam) 
		if !isValid{
			return  utils.ErrorResponse(c, http.StatusBadRequest, "Parameters Validation Failed", http.StatusText(400))
		}
	}

	userRow, err := h.DB.CreateUser(c.Request().Context(), database.CreateUserParams{
		Name: *params.Name,
		Email: *params.Email,
		PhoneNumber: *params.PhoneNumber,
		BusinessName: *params.BusinessName,
		IsAdmin: *params.IsAdmin,
	})

	if err != nil{
		statusCode, msg := utils.ParseDbError(err)
		h.Logger.ErrorContext(c.Request().Context(), "[ERROR]: An error occured:", "db_error", msg, "status_code", statusCode, "email_attempted", *params.Email)
		return  utils.ErrorResponse(c, http.StatusInternalServerError, "An internal Server Error occued...Please try again", "Internal Server Error")
	}

	h.Logger.Info("[SUCCESS] RegisterUser: User profile created...", "user_data", userRow )
	return  utils.SuccessResponse(c, 200, "User Profile Created Successfully", userRow, nil)
}

func(h *Handler) VerifyUser(c echo.Context)error{
	
}