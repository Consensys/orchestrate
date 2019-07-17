package tessera

import (
	"encoding/base64"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

type EnclaveClient struct {
	mux                *sync.Mutex
	endpointForChainID map[string]EnclaveEndpoint
}

type StoreRawResponse struct {
	Key string `json:"key"`
}

func NewEnclaveClient() *EnclaveClient {
	return &EnclaveClient{
		mux:                &sync.Mutex{},
		endpointForChainID: make(map[string]EnclaveEndpoint),
	}
}

func (tc *EnclaveClient) AddClient(chainID string, enclaveClient EnclaveEndpoint) {
	tc.mux.Lock()
	tc.endpointForChainID[chainID] = enclaveClient
	tc.mux.Unlock()
}

func (tc *EnclaveClient) getClient(chainID string) EnclaveEndpoint {
	c, ok := tc.endpointForChainID[chainID]
	if ok {
		return c
	}

	return nil
}

func (tc *EnclaveClient) StoreRaw(chainID string, rawTx []byte, privateFrom string) (txHash []byte, err error) {
	request := map[string]string{
		"payload": base64.StdEncoding.EncodeToString(rawTx),
		"from":    privateFrom,
	}

	c := tc.getClient(chainID)
	if c == nil {
		return nil, fmt.Errorf("no Tessera endpoint for chain id: %s", chainID)
	}

	log.Info("Sending transaction body to 'storeraw' endpoint")

	storeRawResponse := StoreRawResponse{}
	err = c.PostRequest("storeraw", request, &storeRawResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to send a request to Tessera enclave: %s", err)
	}

	txHash, err = base64.StdEncoding.DecodeString(storeRawResponse.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 encoded string in the 'storeraw' response: %s", err)
	}

	return txHash, nil
}

func (tc *EnclaveClient) GetStatus(chainID string) (status string, err error) {
	c := tc.getClient(chainID)
	if c == nil {
		return "", fmt.Errorf("no Tessera endpoint for chain id: %s", chainID)
	}

	log.Infof("Getting Tessera status for the %s chain", chainID)

	return c.GetRequest("upcheck")
}
