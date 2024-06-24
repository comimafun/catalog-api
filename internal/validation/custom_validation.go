package validation

import "github.com/go-playground/validator/v10"

type CustomValidation struct {
	validator *validator.Validate
}

func (cv *CustomValidation) Init() {
	cv.validator.RegisterValidation("url_or_empty", func(fl validator.FieldLevel) bool {
		urlString := fl.Field().String()
		if urlString == "" {
			return true
		}
		return cv.validator.Var(urlString, "url") == nil
	})

	cv.validator.RegisterValidation("day_or_empty", func(fl validator.FieldLevel) bool {
		day := fl.Field().String()
		if day == "" {
			return true
		}
		// day is enum oneof first second both
		return day == "first" || day == "second" || day == "both"
	})
}

func NewCustomValidation(validator *validator.Validate) *CustomValidation {
	return &CustomValidation{validator: validator}
}
