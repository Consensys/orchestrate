package client

import (
	"github.com/ConsenSys/orchestrate/pkg/utils"
)

func (c *HTTPClient) ChainProxyURL(uuid string) string {
	return utils.GetProxyURL(c.config.URL, uuid)
}
