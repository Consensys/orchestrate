package formatters

import (
	"net/http"
	"strings"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

func FormatCreateAccountRequest(req *api.CreateAccountRequest) *entities.Account {
	return &entities.Account{
		Alias:      req.Alias,
		Attributes: req.Attributes,
	}
}

func FormatImportAccountRequest(req *api.ImportAccountRequest) *entities.Account {
	return &entities.Account{
		Alias:      req.Alias,
		Attributes: req.Attributes,
	}
}

func FormatUpdateAccountRequest(req *api.UpdateAccountRequest) *entities.Account {
	return &entities.Account{
		Alias:      req.Alias,
		Attributes: req.Attributes,
	}
}

func FormatAccountResponse(iden *entities.Account) *api.AccountResponse {
	return &api.AccountResponse{
		Alias:               iden.Alias,
		Attributes:          iden.Attributes,
		Address:             iden.Address,
		PublicKey:           iden.PublicKey,
		CompressedPublicKey: iden.CompressedPublicKey,
		TenantID:            iden.TenantID,
		CreatedAt:           iden.CreatedAt,
		UpdatedAt:           iden.UpdatedAt,
	}
}

func FormatAccountFilterRequest(req *http.Request) (*entities.AccountFilters, error) {
	filters := &entities.AccountFilters{}

	qAliases := req.URL.Query().Get("aliases")
	if qAliases != "" {
		filters.Aliases = strings.Split(qAliases, ",")
	}

	if err := utils.GetValidator().Struct(filters); err != nil {
		return nil, err
	}

	return filters, nil
}
