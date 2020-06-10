// +build unit

package multitenancy

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyBuilder_BuildKey(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct { //nolint:maligned // reason
		name         string
		multitenancy bool
		args         args
		want         string
		wantErr      bool
	}{
		{
			"nominal case multitenancy off",
			false,
			args{
				ctx: context.Background(),
				key: "56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E",
			},
			"_56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E",
			false,
		},
		{
			"nominal case multitenancy on",
			true,
			args{
				ctx: WithTenantID(context.Background(), "b49ee1bc-f0fa-430d-89b2-a4fd0dc98906"),
				key: "56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E",
			},
			"b49ee1bc-f0fa-430d-89b2-a4fd0dc9890656202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E",
			false,
		},
		{
			"error case with tenant enable",
			true,
			args{
				ctx: context.Background(),
				key: "56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E",
			},
			"_56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KeyBuilder{
				multitenancy: tt.multitenancy,
			}
			got, err := k.BuildKey(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuildKey() got = %v, want %v", got, tt.want)
			}
			if tt.wantErr {
				if err.Error() != "DB200@: not able to retrieve the tenant UUID: The tenant_id is not present in the Context" {
					t.Errorf("SplitTenant() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}

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
