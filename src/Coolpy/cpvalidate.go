package Coolpy

import (
	"gopkg.in/go-playground/validator.v9"
)

var CpValidate *validator.Validate

func init() {
	CpValidate = validator.New()
}
