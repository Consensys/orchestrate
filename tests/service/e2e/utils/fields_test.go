// +build unit

package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"encoding/json"
)


func TestGetField(t *testing.T) {
	body := `[{"field1":"val1","field2":{"field2.2":"val2"}}, {"field3": "val3"}]` 
	var resp interface{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Error(t)
		return
	}

	actualValue, err := GetField("0.field1", reflect.ValueOf(resp))
	if err != nil {
		t.Error(t)
		return
	}
	
	assert.Equal(t, "val1", actualValue.String())
}
