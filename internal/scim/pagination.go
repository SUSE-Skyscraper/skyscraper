package scim

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type paginationParams struct {
	Offset int32
	Limit  int32
}

var defaultPaginationParams = paginationParams{
	Offset: 0,
	Limit:  10,
}

func paginate(r *http.Request) paginationParams {
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
			return defaultPaginationParams
		}
	}

	if count == "" {
		offset = 0
	} else {
		offset, err = strconv.ParseInt(startIndex, 10, 32)
		if err != nil {
			return defaultPaginationParams
		}
	}

	params := paginationParams{
		Offset: int32(offset),
		Limit:  int32(limit),
	}

	return params
}
