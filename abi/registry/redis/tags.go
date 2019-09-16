package redis

const tagsPrefix = "TagsPrefix"

// TagModel is a zero object gathering methods to look up a abis in redis
type TagModel struct{}

// Tags returns is sugar to manage tags
var Tags = &TagModel{}

// Key serializes a lookup key for a list of tags stored on redis
func (*TagModel) Key(name string) []byte {
	prefixBytes := []byte(tagsPrefix)
	// Allocate memory to build the key
	res := make([]byte, 0, len(prefixBytes)+len(name))
	res = append(res, prefixBytes...)
	res = append(res, name...)
	return res
}

// Get returns a list of tags for a given contract name
func (t *TagModel) Get(conn *Conn, name string) (tags []string, ok bool, err error) {
	tagsBytes, ok, err := conn.LRange(t.Key(name))
	if !ok || err != nil {
		return []string{}, false, err
	}

	// TODO: Make this block error-free. Not all []byte are valid string
	tags = make([]string, len(tagsBytes))
	for index, tagBytes := range tagsBytes {
		tags[index] = string(tagBytes)
	}

	return tags, ok, err
}

// PushIfNotExist stores a new tag in the registry. The function is idemnpotent
func (t *TagModel) PushIfNotExist(conn *Conn, name, tag string) error {
	registeredTags, ok, err := t.Get(conn, name)
	if err != nil {
		return err
	}

	// If a list of tags exists for the tag. Check if the new tag is already registered.
	if ok {
		for _, registeredTag := range registeredTags {
			if registeredTag == tag {
				return nil
			}
		}
	}

	return conn.RPush(t.Key(name), []byte(tag))
}
