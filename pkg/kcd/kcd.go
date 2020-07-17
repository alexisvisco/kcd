package kcd

import (
	"github.com/expectedsh/kcd/pkg/extractor"
	"github.com/expectedsh/kcd/pkg/hook"
)

type H map[string]interface{}

var Config = &struct {
	QueryTag        string
	PathTag         string
	HeaderTag       string
	DefaultTag      string
	QueryExtractor  extractor.Extractor
	PathExtractor   extractor.Extractor
	HeaderExtractor extractor.Extractor
	ErrorHook       hook.Error
	RenderHook      hook.Render
	BindingHook     hook.Bind
	ValidateHook    hook.Validate
}{
	QueryTag:   "query",
	PathTag:    "path",
	HeaderTag:  "header",
	DefaultTag: "default",

	QueryExtractor:  extractor.DefaultQueryExtractor,
	PathExtractor:   extractor.DefaultPathExtractor,
	HeaderExtractor: extractor.DefaultHeaderExtractor,

	ErrorHook:    hook.DefaultErrorHook,
	RenderHook:   hook.DefaultRenderHook,
	BindingHook:  hook.DefaultBindingHook(256 * 1024),
	ValidateHook: hook.DefaultValidateHook,
}
