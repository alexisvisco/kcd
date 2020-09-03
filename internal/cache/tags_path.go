package cache

// TagsPath is a simple map that register from a tag the path of the field.
type TagsPath map[string]string

// Add will add or create the path for a given tag.
// If the tag exist it will add it with the dot notation.
func (t TagsPath) Add(tag string, key string) {
	k, ok := t[tag]
	if !ok {
		k = ""
	}

	if len(k) == 0 {
		k = key
	} else {
		// todo: add a config to modify the dot splitter
		k += "." + key
	}

	t[tag] = k
}

// clone will duplicate this map.
func (t TagsPath) clone() TagsPath {
	n := make(TagsPath, len(t))
	for k, v := range t {
		n[k] = v
	}

	return n
}

// hasValueTag check if there is a key from the set list.
func (t TagsPath) hasValueTag(set []string) bool {
	for _, v := range set {
		_, ok := t[v]
		if ok {
			return true
		}
	}
	return false
}
