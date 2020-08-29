package cache

type TagsPath map[string]string

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

func (t TagsPath) clone() TagsPath {
	n := make(TagsPath, len(t))
	for k, v := range t {
		n[k] = v
	}

	return n
}
