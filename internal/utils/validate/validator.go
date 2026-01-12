package validate

import (
	"sync"

	v "github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

// Same package
var (
	once     sync.Once
	validate *v.Validate
)

func Get() *v.Validate {
	once.Do(func() {
		validate = v.New(v.WithRequiredStructEnabled())
		// register other validations here...
		if err := validate.RegisterValidation("notblank", validators.NotBlank); err != nil {
			panic(err)
		}
	})
	return validate
}
