package httpcache

import (
	"bytes"
	"encoding/gob"
	"net/http"
)

// Response is the cached response data structure.
type response struct {
	Value      []byte
	Header     http.Header
	StatusCode int
}

func newResponse(v []byte, header http.Header, status int) response {
	return response{
		Value:      v,
		Header:     header,
		StatusCode: status,
	}
}

func bytesToResponse(b []byte) (response, error) {
	var r response
	dec := gob.NewDecoder(bytes.NewReader(b))
	err := dec.Decode(&r)
	return r, err
}

// Bytes converts Response data structure into bytes array.
func (r response) toBytes() ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(&r)
	return b.Bytes(), err
}
