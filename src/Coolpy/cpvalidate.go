package Coolpy

import "gopkg.in/go-playground/validator.v8"

var CpValidate *validator.Validate

func init() {
	config := &validator.Config{TagName: "validate"}
	CpValidate = validator.New(config)
}
