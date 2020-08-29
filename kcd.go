package kcd

import (
	"github.com/expectedsh/kcd/pkg/extractor"
	"github.com/expectedsh/kcd/pkg/hook"
)

// Configuration is the main configuration type of kcd
type Configuration struct {
	StringsExtractors []extractor.Strings
	ValueExtractors   []extractor.Value

	ErrorHook    hook.ErrorHook
	BindHook     hook.BindHook
	ValidateHook hook.ValidateHook
	RenderHook   hook.RenderHook
}

// Config is the instance of Configuration type.
// You can add as many extractor you want, modify them ...
// You can set your custom hook too.
var Config = Configuration{
	StringsExtractors: []extractor.Strings{extractor.Path{}, extractor.Header{}, extractor.Query{}},
	ValueExtractors:   []extractor.Value{extractor.Context{}},

	ErrorHook:    hook.Error,
	RenderHook:   hook.Render,
	BindHook:     hook.Bind(256 * 1024),
	ValidateHook: hook.Validate,
}

func (c Configuration) tags() []string {
	tags := make([]string, 0, len(c.StringsExtractors)+len(c.ValueExtractors))

	for _, se := range c.StringsExtractors {
		tags = append(tags, se.Tag())
	}

	for _, ve := range c.ValueExtractors {
		tags = append(tags, ve.Tag())
	}

	return tags
}
