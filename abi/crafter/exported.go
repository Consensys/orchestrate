package crafter

import (
	"sync"
)

var (
	crafter  Crafter
	initOnce = &sync.Once{}
)

// Init initialize ABI Registry
func Init() {
	initOnce.Do(func() {
		// Create crafter
		crafter = &PayloadCrafter{}
	})
}

// SetGlobalCrafter sets global crafter
func SetGlobalCrafter(c Crafter) {
	crafter = c
}

// GlobalCrafter returns global ABI registry
func GlobalCrafter() Crafter {
	return crafter
}
