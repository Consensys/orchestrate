package httputil

import (
	"net/http"
)

func CombineResponseModifiers(modifiers ...func(*http.Response) error) func(*http.Response) error {
	if len(modifiers) > 0 {
		return func(resp *http.Response) error {
			for i := len(modifiers); i > 0; i-- {
				if modifiers[i-1] != nil {
					err := modifiers[i-1](resp)
					if err != nil {
						return err
					}
				}
			}
			return nil
		}
	}

	return func(response *http.Response) error { return nil }
}
