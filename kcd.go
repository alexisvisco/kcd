package kcd

import (
	"github.com/expectedsh/kcd/pkg/extractor"
	"github.com/expectedsh/kcd/pkg/hook"
)

type Configuration struct {
	StringsExtractors []extractor.Strings
	ValueExtractors   []extractor.Value

	ErrorHook    hook.ErrorHook
	BindHook     hook.BindHook
	ValidateHook hook.ValidateHook
	RenderHook   hook.RenderHook
}

var Config = Configuration{
	StringsExtractors: []extractor.Strings{extractor.Path{}, extractor.Header{}, extractor.Query{}},
	ValueExtractors:   []extractor.Value{extractor.Context{}},

	ErrorHook:    hook.Error,
	RenderHook:   hook.Render,
	BindHook:     hook.Bind(256 * 1024),
	ValidateHook: hook.Validate,
}
