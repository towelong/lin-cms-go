package validator

import (
	"github.com/gin-gonic/gin/binding"
	cn "github.com/go-playground/locales/zh_Hans_CN"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"sync"
)

var (
	Trans ut.Translator
)

type DefaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &DefaultValidator{}

func InitValidator() {
	binding.Validator = New()
}

func New() *DefaultValidator {
	defaultValidator := new(DefaultValidator)
	defaultValidator.lazyInit()
	translator := cn.New()
	uni := ut.New(translator)
	Trans, _ = uni.GetTranslator("zh_Hans_CN")
	_ = zhTranslations.RegisterDefaultTranslations(defaultValidator.validate, Trans)
	return defaultValidator
}


func (v *DefaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyInit()
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

func (v *DefaultValidator) Engine() interface{} {
	v.lazyInit()
	return v.validate
}

func (v *DefaultValidator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("validate")
		// add any custom validations etc. here

		v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("label")
			if name == "" {
				return "字段"
			}
			return name
		})
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
