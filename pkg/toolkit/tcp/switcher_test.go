// +build unit

package tcp_test

import (
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/tcp"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/tcp/mock"
	"github.com/golang/mock/gomock"
)

func TestSwitcher(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	switcher := tcp.NewSwitcher()
	conn := mock.NewMockWriteCloser(ctrlr)
	conn.EXPECT().Close().Times(1)
	switcher.ServeTCP(conn)

	h1 := mock.NewMockHandler(ctrlr)
	switcher.Switch(h1)
	h1.EXPECT().ServeTCP(conn)
	switcher.ServeTCP(conn)

	h2 := mock.NewMockHandler(ctrlr)
	switcher.Switch(h2)
	h2.EXPECT().ServeTCP(conn)
	switcher.ServeTCP(conn)
}
