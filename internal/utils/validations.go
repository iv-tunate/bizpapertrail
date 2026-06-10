package utils

import (
	"log"
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
	log.Println("Registering validators")
	
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

func ValidateUserParams(param any)(bool){
	err := validate.Struct(param);
	if err != nil{
		return  false
	}
	return  true
}