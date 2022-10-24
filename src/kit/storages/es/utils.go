package es

import (
	"github.com/africarealty/server/src/kit"
	"github.com/olivere/elastic/v7"
)

var Missings = map[string]interface{}{
	kit.SortRequestMissingFirst: EsSortRequestMissingFirst,
	kit.SortRequestMissingLast:  EsSortRequestMissingLast,
	"":                          nil,
}

func ToSortRequestEs(request *kit.SortRequest) (*elastic.SortInfo, error) {
	if request.Field == "" {
		return nil, ErrEsSortRequestFieldEmpty()
	}

	res, ok := Missings[request.Missing]
	if !ok {
		return nil, ErrEsSortRequestMissingInvalid(request.Missing)
	}

	return &elastic.SortInfo{
		Field:     request.Field,
		Ascending: request.Asc,
		Missing:   res,
	}, nil
}
