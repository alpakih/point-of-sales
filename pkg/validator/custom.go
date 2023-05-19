package validator

import (
	"regexp"
	"strings"
	"time"

	validatorGo "github.com/go-playground/validator/v10"
)

func registerCustomValidator(v *validatorGo.Validate) {

	if err := v.RegisterValidation("email_address", ValidateEmail); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("number_format", ValidateOnlyNumber); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("mobile_phone", ValidateMobilePhone); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("date_only", ValidateDateOnly); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("must_lte_current_date", MustLteCurrentDate); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("lte_current_date", LteCurrentDate); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("only_words", ValidateStringOnlyWords); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("no_space", ValidateNoSpace); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("name", ValidateName); err != nil {
		panic(err)
	}
}

func ValidateName(fl validatorGo.FieldLevel) bool {
	if fl.Field().String() != "" {
		regex := regexp.MustCompile(`^[a-zA-Z\s,.’' ]*$`)
		return regex.MatchString(fl.Field().String())
	}
	return true
}

func ValidateEmail(fl validatorGo.FieldLevel) bool {
	if fl.Field().String() != "" {
		regex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+")
		return regex.MatchString(fl.Field().String())
	}
	return true
}

func ValidateMobilePhone(fl validatorGo.FieldLevel) bool {
	if fl.Field().String() != "" {
		regex := regexp.MustCompile(`^08\d{7,12}$`)
		return regex.MatchString(fl.Field().String())
	}
	return true
}

func ValidateDateOnly(fl validatorGo.FieldLevel) bool {
	if fl.Field().String() != "" {
		parse, err := time.Parse("2006-01-02", fl.Field().String())
		if err != nil {
			return false
		}
		if parse.IsZero() {
			return false
		}
		regex := regexp.MustCompile(`^\d{4}-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])$`)
		return regex.MatchString(fl.Field().String())
	}
	return true
}

func MustLteCurrentDate(fl validatorGo.FieldLevel) bool {
	if fl.Field().String() != "" {

		if parse, err := time.Parse("2006-01-02", fl.Field().String()); err != nil {
			return false
		} else {
			if parse.IsZero() {
				return false
			}

			currentDate, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))

			if err != nil {
				return false
			}
			if currentDate.IsZero() {
				return false
			}

			if parse.Sub(currentDate) > 0 {
				return false
			}
		}
	}
	return true
}

func LteCurrentDate(fl validatorGo.FieldLevel) bool {

	if fl.Field().String() != "" {

		if parse, err := time.Parse("2006-01-02", fl.Field().String()); err != nil {
			return false
		} else {
			if parse.After(time.Now()) {
				return false
			}
			return true
		}
	}
	return true
}

func ValidateOnlyNumber(fl validatorGo.FieldLevel) bool {
	if fl.Field().String() != "" {
		regex := regexp.MustCompile(`^[0-9]*$`)
		return regex.MatchString(fl.Field().String())
	}
	return true
}

func ValidateNoSpace(field validatorGo.FieldLevel) bool {
	value := field.Field().String()

	res := strings.TrimSpace(value)
	if len(res) == 0 {
		return false
	}

	return true
}

func ValidateStringOnlyWords(fl validatorGo.FieldLevel) bool {
	if fl.Field().String() != "" {
		regex := regexp.MustCompile(`^[A-Za-z ]+$`)
		return regex.MatchString(fl.Field().String())
	}
	return true

}
