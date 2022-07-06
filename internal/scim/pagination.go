package scim

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type paginationParams struct {
	Offset int32
	Limit  int32
}

func paginate(r *http.Request) (paginationParams, error) {
	startIndex := chi.URLParam(r, "startIndex")
	count := chi.URLParam(r, "count")

	var err error
	var limit int64
	var offset int64

	if startIndex == "" {
		limit = 10
	} else {
		limit, err = strconv.ParseInt(count, 10, 32)
		if err != nil {
			return paginationParams{}, errors.New("invalid count")
		}
	}

	if count == "" {
		offset = 0
	} else {
		offset, err = strconv.ParseInt(startIndex, 10, 32)
		if err != nil {
			return paginationParams{}, errors.New("invalid startIndex")
		}
	}

	params := paginationParams{
		Offset: int32(offset),
		Limit:  int32(limit),
	}

	return params, nil
}
