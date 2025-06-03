package validators

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("gender", validateGender)
	v.RegisterValidation("transaction_type", validateTransactionType)
	v.RegisterValidation("date", DateValidator)
}

func validateGender(fl validator.FieldLevel) bool {
	gender := fl.Field().String()
	return gender == "L" || gender == "P"
}

func validateTransactionType(fl validator.FieldLevel) bool {
	transactionType := fl.Field().String()
	return transactionType == "stock_in" || transactionType == "stock_out"
}

func DateValidator(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}
