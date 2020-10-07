package testutils

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store/models"

	"github.com/gofrs/uuid"
)

func FakeIdentityModel() *models.Identity {
	return &models.Identity{
		UUID: uuid.Must(uuid.NewV4()).String(),
	}
}
