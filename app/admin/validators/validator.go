package validators

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

func CheckCtx(ctx *gin.Context, s interface{}) error {
	err := Validate.StructCtx(ctx, s)
	if err != nil {
		return err
	}
	return nil
}
