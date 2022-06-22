package utils

import (
	"net/http"
	"strconv"

	"github.com/consensys/orchestrate/pkg/types/entities"

	"github.com/consensys/orchestrate/pkg/errors"
)

func FilterIntegerValueWithKey(req *http.Request) (*entities.PaginationFilters, error) {
	filters := &entities.PaginationFilters{}
	var err error
	qPage := req.URL.Query().Get("page")
	if qPage != "" {
		filters.Page, err = strconv.Atoi(qPage)
		if err != nil || filters.Page < 0 {
			return filters, errors.InvalidFormatError("page format is invalid, must be positive integer")
		}
	}
	qLimit := req.URL.Query().Get("limit")
	if qLimit != "" {
		filters.Limit, err = strconv.Atoi(qLimit)
		if err != nil || filters.Limit < 0 {
			return filters, errors.InvalidFormatError("limit format is invalid, must be positive integer")
		}
	}
	return filters, nil
}
