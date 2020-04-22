package utils

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// func AddFlagOnce(flgs *pflag.FlagSet, name string, dfault interface{}, desc string, viperKey string) {
// 	if flgs.Lookup(name) == nil {
// 		switch val := dfault.(type) {
// 		case nil:
// 			panic(fmt.Sprint("invalid type of %v", val))
// 		case int:
// 			flgs.Int(name, val, desc)
// 		case uint:
// 			flgs.Uint(name, val, desc)
// 		case string:
// 			flgs.String(name, val, desc)
// 		case []string:
// 			flgs.StringSlice(name, val, desc)
// 		default:
// 			panic(fmt.Sprint("invalid type of %v", val)) // here v has type interface{}
// 		}
// 	}
// 
// 	_ = viper.BindPFlag(viperKey, flgs.Lookup(name))
// }

func PreRunBindFlags(vipr *viper.Viper, flgs *pflag.FlagSet, ignore string) {
	for _, vk := range vipr.AllKeys() {
		if ignore != "" && strings.HasPrefix(vk, ignore) {
			continue
		}

		// Convert viperKey to cmd flag name
		// For example: 'rest.api' to "rest-api"
		name := strings.Replace(vk, ".", "-", -1)
		
		// Only bind in case command flags contain the name
		if flgs.Lookup(name) != nil {
			_ = viper.BindPFlag(vk, flgs.Lookup(name))
		}
	}
}
