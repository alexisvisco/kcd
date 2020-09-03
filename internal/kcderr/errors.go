package kcderr

import (
	"github.com/expectedsh/errors"
)

var (
	// Input errors are related to bad request and didn't need to be logged
	Input = errors.Kind("input")

	// InputCritical should never happen, those kind of error are logged
	InputCritical = errors.Kind("input_critical")

	// OutputCritical should never happen, those kind of error are logged
	OutputCritical = errors.Kind("output_critical")
)
