package swagger

import (
	"sync"
)

var (
	handler  *Handler
	initOnce = &sync.Once{}
)

// Initialize API handlers
func Init(specsPath, uiPath string) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Set Swagger handler
		handler = NewHandler(specsPath, uiPath)
	})
}

// GlobalHandler return the swagger
func GlobalHandler() *Handler {
	return handler
}
