package redis

const catalogPrefix = "catalogPrefix"

// CatalogModel is a zero object gathering methods to look up a abis in redis
type CatalogModel struct{}

// Catalog returns is sugar to return a CatalogModel object
var Catalog = &CatalogModel{}

// Key serializes a lookup key for the contract name catalog stored on redis
func (*CatalogModel) Key() []byte {
	prefixBytes := []byte(catalogPrefix)
	// Allocate memory to build the key
	res := make([]byte, 0, len(prefixBytes))
	res = append(res, prefixBytes...)
	return res
}

// Get returns a list of contract name
func (t *CatalogModel) Get(conn *Conn) (names []string, ok bool, err error) {
	namesBytes, ok, err := conn.LRange(t.Key())
	if !ok || err != nil {
		return []string{}, false, err
	}

	// TODO: Make this block error-free. Not all []byte are valid string
	names = make([]string, len(namesBytes))
	for index, nameBytes := range namesBytes {
		names[index] = string(nameBytes)
	}

	return names, ok, err
}

// PushIfNotExist push a new contract name in the registry. The function is idemnpotent
func (t *CatalogModel) PushIfNotExist(conn *Conn, name string) error {
	registeredTags, ok, err := t.Get(conn)
	if err != nil {
		return err
	}

	// If a list of tags exists for the name. Check if the new name is already registered.
	if ok {
		for _, registeredNames := range registeredTags {
			if registeredNames == name {
				return nil
			}
		}
	}

	return conn.LPush(t.Key(), []byte(name))
}
