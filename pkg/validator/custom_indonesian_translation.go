package validator

import (
	"log"
	"strings"

	ut "github.com/go-playground/universal-translator"

	validatorGo "github.com/go-playground/validator/v10"
)

func registerCustomIndonesianTranslator(v *validatorGo.Validate, trans ut.Translator) {

	if err := v.RegisterTranslation("date_only", trans, func(ut ut.Translator) error {
		if err := ut.Add("date_only", "{0} tidak valid atau format tidak sesuai (yyyy-mm-dd).", false); err != nil {
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
		if err := ut.Add("must_lte_current_date", "{0} harus lebih kecil atau sama dengan tanggal sekarang.", false); err != nil {
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
		if err := ut.Add("number_format", "{0} hanya angka yang diperbolehkan.", false); err != nil {
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
		if err := ut.Add("mobile_phone", "{0} format tidak valid.", false); err != nil {
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
		if err := ut.Add("email_address", "{0} format tidak valid.", false); err != nil {
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
		if err := ut.Add("only_words", "{0} tidak diperbolehkan mengandung angka atau tanda baca.", false); err != nil {
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
		if err := ut.Add("no_space", "field tidak bisa hanya berisi spasi.", false); err != nil {
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
		if err := ut.Add("enum", "kriteria yang diterima dari {0} adalah {1}", false); err != nil {
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
		if err := ut.Add("name", "{0} tidak diperbolehkan mengandung angka atau tanda baca selain '.' dan ','.", false); err != nil {
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
