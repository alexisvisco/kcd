package cache

import (
	"encoding/json"
	"reflect"

	"github.com/alexisvisco/kcd/internal/types"
)

// StructAnalyzer is the cache struct analyzer
type StructAnalyzer struct {
	tags           []string
	valueTag       []string
	mainStructType reflect.Type
}

// NewStructAnalyzer instantiate a new StructAnalyzer
func NewStructAnalyzer(stringsTags, valueTags []string, mainStructType reflect.Type) *StructAnalyzer {
	return &StructAnalyzer{
		tags:           append(stringsTags, valueTags...),
		valueTag:       valueTags,
		mainStructType: mainStructType,
	}
}

// StructCache
type StructCache struct {
	IsRoot     bool
	Index      []int
	Resolvable []FieldMetadata

	Child []StructCache
}

// String is a way to debug StructCache
func (s StructCache) String() string {
	marshal, _ := json.MarshalIndent(s, "", " ")
	return string(marshal)
}

// FieldMetadata contains all the necessary field to decode
type FieldMetadata struct {
	Index                 []int
	Paths                 TagsPath
	Type                  reflect.Type
	ImplementUnmarshaller bool
	ArrayOrSlice          bool
	DefaultValue          string
	Exploder              string
}

func (f FieldMetadata) GetDefaultFieldName() string {
	for tag, path := range f.Paths {
		if tag != "default" {
			return path
		}
	}
	return "unknown"
}

// Cache will take all fields which contain a tag to be lookup.
func (s StructAnalyzer) Cache() StructCache {
	sc := newStructCache()

	s.cache(&sc, TagsPath{}, s.mainStructType)

	return sc
}

func (s StructAnalyzer) cache(cache *StructCache, paths TagsPath, t reflect.Type) (containTags bool) {
	if t == nil {
		return false
	}

	if !sanitizePtrType(&t) {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		var (
			structField = t.Field(i)
			metadata    = FieldMetadata{Index: structField.Index, Type: structField.Type}
		)

		if types.IsImplementingUnmarshaler(metadata.Type) {
			metadata.ImplementUnmarshaller = true
		}

		if !sanitizePtrType(&metadata.Type) {
			continue
		}

		fieldHasTag := false
		currentPaths := paths.clone()

		if s.lookupTags(structField, currentPaths) {
			fieldHasTag = true
			containTags = true
		}

		hasValueTag := currentPaths.hasValueTag(s.valueTag)

		if !hasValueTag && (structField.Anonymous || metadata.Type.Kind() == reflect.Struct) {
			childStructCache := newStructCacheFromField(structField)
			childStructContainTag := s.cache(&childStructCache, currentPaths, metadata.Type)

			if childStructContainTag {
				cache.Child = append(cache.Child, childStructCache)
				fieldHasTag = true
				containTags = true
			}

			if childStructContainTag || !metadata.ImplementUnmarshaller {
				continue
			}
		}

		if !metadata.ImplementUnmarshaller &&
			(metadata.Type.Kind() == reflect.Slice || metadata.Type.Kind() == reflect.Array) {
			typeOfArray := metadata.Type.Elem()
			if !sanitizePtrType(&typeOfArray) {
				continue
			}

			metadata.ArrayOrSlice = true
			metadata.Type = typeOfArray
		}

		if !(hasValueTag || (fieldHasTag && (metadata.ImplementUnmarshaller || types.IsUnmarshallable(metadata.Type)))) {
			continue
		}

		metadata.DefaultValue = structField.Tag.Get("default")
		metadata.Exploder = structField.Tag.Get("exploder")
		metadata.Paths = currentPaths

		cache.Resolvable = append(cache.Resolvable, metadata)
	}

	return containTags
}

func (s StructAnalyzer) lookupTags(structField reflect.StructField, currentPaths TagsPath) (containTags bool) {
	hasTags := false
	for _, tag := range s.tags {
		lookup, ok := structField.Tag.Lookup(tag)
		if ok {
			hasTags = true
			currentPaths.Add(tag, lookup)
		}
	}
	return hasTags
}

func newStructCacheFromField(field reflect.StructField) StructCache {
	return StructCache{
		Index:      field.Index,
		Resolvable: []FieldMetadata{},
		Child:      []StructCache{},
	}
}

func newStructCache() StructCache {
	return StructCache{
		IsRoot:     true,
		Resolvable: []FieldMetadata{},
		Child:      []StructCache{},
	}
}

func sanitizePtrType(t *reflect.Type) (success bool) {
	x := *t
	if x.Kind() == reflect.Ptr {
		*t = x.Elem()

		if (*t).Kind() == reflect.Ptr {
			return false
		}
	}
	return true
}
