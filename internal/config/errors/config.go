package errors

import "errors"

// FileInvalid represents that a TOML parsing error occurred
var FileInvalid = errors.New("config file is invalid")

// FatalValidation represents that a fatal validation error occurred
var FatalValidation = errors.New("fatal validation errors")
