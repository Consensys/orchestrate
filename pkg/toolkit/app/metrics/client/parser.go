package client

import (
	"net/http"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

func ParseResponse(resp *http.Response) (map[string]*dto.MetricFamily, error) {
	var parser expfmt.TextParser
	mFams, err := parser.TextToMetricFamilies(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	famis := map[string]*dto.MetricFamily{}
	for mid, mf := range mFams {
		famis[mid] = mf
	}

	return famis, nil
}
