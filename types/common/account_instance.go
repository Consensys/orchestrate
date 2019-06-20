package common

import "fmt"

// Short returns a string representation of the account instance
func (instance *AccountInstance) Short() (string, error) {
	addr := instance.GetAccount().Address()
	return fmt.Sprintf("%v@%v", addr.String(), instance.GetChain().GetId()), nil
}
