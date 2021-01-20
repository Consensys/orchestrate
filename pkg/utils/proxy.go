package utils

import (
	"fmt"
)

func GetProxyURL(proxyURL, chainUUID string) string {
	return fmt.Sprintf("%s/proxy/chains/%s", proxyURL, chainUUID)
}

func GetProxyTesseraURL(proxyURL, chainUUID string) string {
	return fmt.Sprintf("%s/proxy/chains/tessera/%s", proxyURL, chainUUID)
}
