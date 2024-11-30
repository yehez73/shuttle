package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidateStruct(c *fiber.Ctx, v interface{}) error {
    validate := validator.New()
    if err := validate.Struct(v); err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            return BadRequestResponse(c, err.Field()+" is "+err.Tag(), nil)
        }
    }
    return nil
}