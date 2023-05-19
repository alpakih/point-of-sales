package validator

import (
	"log"
	"strings"

	ut "github.com/go-playground/universal-translator"

	validatorGo "github.com/go-playground/validator/v10"
)

func registerCustomEnglishTranslator(v *validatorGo.Validate, trans ut.Translator) {

	if err := v.RegisterTranslation("date_only", trans, func(ut ut.Translator) error {
		if err := ut.Add("date_only", "{0} invalid or invalid date format (yyyy-mm-dd).", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}

	if err := v.RegisterTranslation("must_lte_current_date", trans, func(ut ut.Translator) error {
		if err := ut.Add("must_lte_current_date", "{0} must lower than or equal current date.", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}

	if err := v.RegisterTranslation("number_format", trans, func(ut ut.Translator) error {
		if err := ut.Add("number_format", "{0} only number accepted.", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}

	if err := v.RegisterTranslation("mobile_phone", trans, func(ut ut.Translator) error {
		if err := ut.Add("mobile_phone", "{0} invalid format mobile phone.", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}

	if err := v.RegisterTranslation("email_address", trans, func(ut ut.Translator) error {
		if err := ut.Add("email_address", "{0} invalid format email.", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}

	if err := v.RegisterTranslation("only_words", trans, func(ut ut.Translator) error {
		if err := ut.Add("only_words", "{0} must not contain numbers or punctuation mark.", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}

	if err := v.RegisterTranslation("no_space", trans, func(ut ut.Translator) error {
		if err := ut.Add("no_space", "field can't be filled with only space.", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}

	if err := v.RegisterTranslation("enum", trans, func(ut ut.Translator) error {
		if err := ut.Add("enum", "acceptance value of {0} is {1}", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		// first, clean/remove the comma
		cleaned := strings.Replace(fe.Param(), "-", " ", -1)

		// convert 'cleaned' comma separated string to slice
		strSlice := strings.Fields(cleaned)

		t, err := ut.T(fe.Tag(), fe.Field(), strings.Join(strSlice, ","))
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}

	if err := v.RegisterTranslation("name", trans, func(ut ut.Translator) error {
		if err := ut.Add("name", "{0} must not contain numbers or punctuation mark other than '.' or ','.", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}

}
