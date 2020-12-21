package utils

import (
	"fmt"
)

func GetProxyURL(chainRegistryURL, chainUUID string) string {
	return fmt.Sprintf("%s/%s", chainRegistryURL, chainUUID)
}

func GetProxyTesseraURL(chainRegistryURL, chainUUID string) string {
	return fmt.Sprintf("%s/tessera/%s", chainRegistryURL, chainUUID)
}
