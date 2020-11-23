// +build unit

package switcher

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/handler/mock"
)

func TestSwitcher(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	switcher := New()
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	rw := httptest.NewRecorder()

	h1 := mock.NewMockHandler(ctrlr)
	switcher.Switch(h1)
	h1.EXPECT().ServeHTTP(rw, req)
	switcher.ServeHTTP(rw, req)

	h2 := mock.NewMockHandler(ctrlr)
	switcher.Switch(h2)
	h2.EXPECT().ServeHTTP(rw, req)
	switcher.ServeHTTP(rw, req)
}
