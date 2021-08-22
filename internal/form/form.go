// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"net/http"

	"github.com/flamego/binding"
	"github.com/flamego/flamego"
	"github.com/flamego/validator"
	en_translations "github.com/flamego/validator/translations/en"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	jsoniter "github.com/json-iterator/go"
	log "unknwon.dev/clog/v2"
)

func Bind(model interface{}) flamego.Handler {
	validate := validator.New()
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	validate = validator.New()
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	return binding.JSON(model, binding.Options{
		ErrorHandler: errorHandler(trans),
		Validator:    validate,
	})
}

func errorHandler(trans ut.Translator) flamego.Handler {
	return func(c flamego.Context, errors binding.Errors) {
		c.ResponseWriter().WriteHeader(http.StatusBadRequest)
		c.ResponseWriter().Header().Set("Content-Type", "application/json")

		var errorCode int
		var msg string
		if errors[0].Category == binding.ErrorCategoryDeserialization {
			errorCode = 40000
			msg = "Error payload"
		} else {
			errorCode = 40001
			errs := errors[0].Err.(validator.ValidationErrors)
			for _, e := range errs {
				msg = e.Translate(trans)
				break
			}
		}

		body := map[string]interface{}{
			"error": errorCode,
			"msg":   msg,
		}
		err := jsoniter.NewEncoder(c.ResponseWriter()).Encode(body)
		if err != nil {
			log.Error("Failed to encode response body: %v", err)
		}
	}
}
