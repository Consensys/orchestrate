package tessera

type Client interface {
	// AddClient adds a Tessera client for a specified chain UUID
	AddClient(chainID string, tesseraEndpoint EnclaveEndpoint)

	// StoreRaw stores "data" field of a transaction in Tessera privacy enclave
	// It returns a hash of a stored transaction that should be used instead of transaction data
	StoreRaw(chainID string, rawTx []byte, privateFrom string) (txHash []byte, err error)

	// GetStatus returns status of Tessera enclave if it is up or an error if it is down
	GetStatus(chainID string) (status string, err error)
}

//go:generate mockgen -destination=mock/mock.go -package=mock . EnclaveEndpoint

type EnclaveEndpoint interface {
	// PostRequest - sends a request to Tessera private enclave
	// path - a URL path in a request to send
	// request - a body of a request to send
	// response - a pointer to an object that will be populated with the parsed body of a response
	PostRequest(path string, request, response interface{}) error

	// GetRequest - sends a GET request to Tessera private enclave
	// path - a URL path in a request to send
	// Returns body of a response as a string and an error if one occurred
	GetRequest(path string) (string, error)
}
