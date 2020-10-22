// +build unit

package multitenancy

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitTenant(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		wantRes1 context.Context
		wantRes2 string
		wantErr  bool
	}{
		{
			"pkey with tenant",
			"56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E@b49ee1bc-f0fa-430d-89b2-a4fd0dc98906",
			WithTenantID(context.Background(), "b49ee1bc-f0fa-430d-89b2-a4fd0dc98906"),
			"56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E",
			false,
		},
		{
			"pkey without tenant",
			"56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E",
			context.Background(),
			"56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E",
			false,
		},
		{
			"error case",
			"b49ee1bc-f0fa-430d-89b2-a4fd0dc98906@b49ee1bc-f0fa-430d-89b2-a4fd0dc98906@56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E",
			nil,
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2, err := SplitTenant(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got1, tt.wantRes1) {
				t.Errorf("SplitTenant() got1 = %v, wantRes1 %v", got1, tt.wantRes1)
			}
			if got2 != tt.wantRes2 {
				t.Errorf("SplitTenant() got2 = %v, wantRes1 %v", got2, tt.wantRes2)
			}
			if tt.wantErr {
				assert.Error(t, err, "SplitTenant() error = %v, wantErr %v")
			}
		})
	}
}
