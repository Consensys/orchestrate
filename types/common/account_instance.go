package common

import "fmt"

// Short returns a string representation of the account instance
func (instance *AccountInstance) Short() string {
	var addr, id string
	if instance.GetAccount() == nil {
		addr = ""
	} else {
		addr = instance.GetAccount().Address().String()
	}
	if instance.GetChain() == nil {
		id = ""
	} else {
		id = string(instance.GetChain().GetChainId())
	}
	return fmt.Sprintf("%v@%v", addr, id)
}
