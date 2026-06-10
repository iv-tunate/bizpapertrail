package utils

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	PasswordRegex = regexp.MustCompile(`[A-Za-z].*[0-9]|[0-9].*[A-Za-z]`)
	PhoneRegex = regexp.MustCompile(`^(\+?[1-9]\d{1,14}|0\d{9,10})$`)
	EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func ValidatePhoneNumber(phone_number *string) (bool){
	
	phone := strings.TrimSpace(*phone_number)
	return PhoneRegex.MatchString(phone)
}


var validate = validator.New();

func init(){
	validate.RegisterValidation("email_regex", func(f1 validator.FieldLevel) bool{
		return EmailRegex.MatchString(f1.Field().String())
	})
	validate.RegisterValidation("password", func(f1 validator.FieldLevel) bool{
		return PasswordRegex.MatchString(f1.Field().String())
	})
	validate.RegisterValidation("phone_regex", func(f1 validator.FieldLevel) bool{
		return PhoneRegex.MatchString(f1.Field().String())
	})
}

func ValidateUserParams(param any) map[string]string {
    err := validate.Struct(param)
    if err == nil {
        return nil
    }

    errors := make(map[string]string)
    for _, e := range err.(validator.ValidationErrors) {
        field := strings.ToLower(e.Field())
        switch e.Tag() {
        case "required":
            errors[field] = field + " is required"
        case "min":
            errors[field] = field + " must be at least " + e.Param() + " characters"
        case "email_regex":
            errors[field] = "invalid email format"
        case "password":
            errors[field] = "password must contain uppercase, lowercase, number and special character"
        case "phone_regex":
            errors[field] = "invalid phone number format"
        default:
            errors[field] = field + " is invalid"
        }
    }
    return errors
}