package util

import (
	"github.com/go-playground/locales/zh_Hans_CN"
	unTrans "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"log"
	"reflect"
)

// Validate 数据校验
func Validate(data interface{}) (s string, err error) {
	validate := validator.New()
	uni := unTrans.New(zh_Hans_CN.New())
	trans, _ := uni.GetTranslator("zh_Hans_CN")

	if err = zhTrans.RegisterDefaultTranslations(validate, trans); err != nil {
		log.Println("validate err:", err)
		return
	}
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		return label
	})

	if err = validate.Struct(data); err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			s = v.Translate(trans)
		}
	}
	return
}
