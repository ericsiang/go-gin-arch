package validlang

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	unitrans "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	chTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var trans unitrans.Translator

// InitValidateLang 初始化驗證語系，locales 通常取决于 http 请求头的 'Accept-Language'
func InitValidateLang(local string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New() //chinese
		enT := en.New() //english
		unitrans := unitrans.New(enT, zhT, enT)

		var ok bool
		trans, ok = unitrans.GetTranslator(local)
		if !ok {
			return fmt.Errorf("unitrans.GetTranslator(%s) failed", local)
		}

		//register translate
		// 注册翻译器
		switch local {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = chTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return

	}
	return
}

// ErrorValidateCheckAndTrans 錯誤驗證檢查並翻譯
func ErrorValidateCheckAndTrans(err error) (translateErrs validator.ValidationErrorsTranslations, ok bool) {
	// 取得validator.ValidationErrors類型的errors，
	validErrs, ok := err.(validator.ValidationErrors)
	if ok { //是validator.ValidationErrors類型錯誤則進行翻譯
		translateErrs = validErrs.Translate(trans)
		return translateErrs, true
	}

	return nil, false
}
