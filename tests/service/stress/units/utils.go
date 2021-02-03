package units

func constructorArgs(contractName string) []interface{} {
	switch contractName {
	case "ERC20", "ERC721":
		return []interface{}{"Name", "Symbol"}
	case "ERC777":
		return []interface{}{"Name", "Symbol", []string{}}
	default:
		return []interface{}{}
	}
}
