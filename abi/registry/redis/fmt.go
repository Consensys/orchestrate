package redis

// Prefixes are there to prevent between several mapping
// if they happens to have identical keys layout
const methodKeyPrefix = "MethodKey"
const eventKeyPrefix = "EventKey"
const codeHashPrefix = "CodeHash"

func methodKey(codeHash []byte, selector [4]byte) []byte {
	// Allocate memory to build the key
	res := make([]byte, 0,
		len(methodKeyPrefix)+len(codeHash)+4)
	// Fill in the buffer and return it
	res = append(res, methodKeyPrefix...)
	res = append(res, codeHash...)
	return append(res, selector[:]...)
}

func eventKey(codeHash, eventID []byte, index uint) []byte {
	indexBytes := []byte{byte(index)}
	// Allocate memory to build the key
	res := make([]byte, 0,
		len(eventKeyPrefix)+len(codeHash)+len(eventID)+len(indexBytes))
	// Fill in the buffer and return it
	res = append(res, eventKeyPrefix...)
	res = append(res, codeHash...)
	res = append(res, eventID...)
	return append(res, indexBytes...)
}

func codeHashKey(chainID string, address []byte) []byte {
	chainIDBytes := []byte(chainID)
	// Allocate memory to build the key
	res := make([]byte, 0,
		len(codeHashPrefix)+len(chainIDBytes)+len(address))
	// Fill in the buffer and returns it
	res = append(res, codeHashPrefix...)
	res = append(res, chainIDBytes...)
	res = append(res, address...)
	return res
}
