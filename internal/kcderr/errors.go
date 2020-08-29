package kcderr

import (
	"github.com/expectedsh/errors"
)

var (
	Input          = errors.Kind("input")
	InputCritical  = errors.Kind("input_critical")
	OutputCritical = errors.Kind("output_critical")
)
