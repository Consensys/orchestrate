package common

import (
	"fmt"
)

// StringShort returns a string representation of the method
func (method *Method) StringShort() string {
	if method.GetContract() == "" {
		return ""
	}

	if method.GetVersion() == "" {
		if method.GetDeploy() {
			return fmt.Sprintf("deploy(%v)", method.GetContract())
		}
		return fmt.Sprintf("%v@%v", method.GetName(), method.GetContract())
	}

	if method.GetDeploy() {
		return fmt.Sprintf("deploy(%v[%v])", method.GetContract(), method.GetVersion())
	}

	return fmt.Sprintf("%v@%v[%v]", method.GetName(), method.GetContract(), method.GetVersion())
}
