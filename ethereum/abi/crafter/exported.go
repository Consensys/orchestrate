package crafter

import (
	"sync"
)

const component = "abi.crafter"

var (
	crafter  Crafter
	initOnce = &sync.Once{}
)

// Init initialize Crafter
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

// GlobalCrafter returns global handler
func GlobalCrafter() Crafter {
	return crafter
}
