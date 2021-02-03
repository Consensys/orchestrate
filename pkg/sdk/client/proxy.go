package client

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

func (c *HTTPClient) ChainProxyURL(uuid string) string {
	return utils.GetProxyURL(c.config.URL, uuid)
}
