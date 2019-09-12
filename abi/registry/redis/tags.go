package redis

const tagsPrefix = "TagsPrefix"

// TagModel is a zero object gathering methods to look up a abis in redis
type TagModel struct{}

// Tag returns is sugar to return an abi object
var Tag = &TagModel{}

// Key serializes a lookup key for an ABI stored on redis
func (*TagModel) Key(name string) []byte {
	prefixBytes := []byte(tagsPrefix)
	// Allocate memory to build the key
	res := make([]byte, 0, len(prefixBytes)+len(name))
	res = append(res, prefixBytes...)
	res = append(res, name...)
	return res
}

// Get returns a serialized contract from its corresponding bytecode hash
func (t *TagModel) Get(conn *Conn, name string) ([][]byte, error) {
	return conn.LRange(t.Key(name))
}

// Add stores an abi object in the registry
func (t *TagModel) Add(conn *Conn, name string, tagBytes []byte) error {
	return conn.LPush(t.Key(name), tagBytes)
}
