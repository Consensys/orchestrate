package ethclient

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
)

// Define a service that will be used as mock
type Service struct {
	Raw string
}

func (s *Service) SendRawTransaction(ctx context.Context, raw string) {
	s.Raw = raw
}

func NewTestServer(name string, service interface{}) (*rpc.Server, error) {
	server := rpc.NewServer()

	err := server.RegisterName(name, service)
	if err != nil {
		return nil, err
	}

	return server, nil
}
func TestClient(t *testing.T) {

	// Create a test server and register mock service
	service := new(Service)
	server, err := NewTestServer("eth", service)

	if err != nil {
		t.Errorf("NewServer: %v", err)
	}

	// Create a rpc client connected to our test server
	c := rpc.DialInProc(server)
	client := NewClient(c)

	// Send a raw transaction
	raw := "0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80"
	err = client.SendRawTransaction(context.Background(), raw)
	if err != nil {
		t.Errorf("SendRawTransaction: %v", err)
	}

	// Check service has been correctly updated
	if service.Raw != raw {
		t.Errorf("SendRawTransaction: expected %q but got %q", raw, service.Raw)
	}
}
