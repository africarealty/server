package es

import (
	"github.com/africarealty/server/src/kit"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ToSortRequestEs(t *testing.T) {
	tests := []struct {
		name string
		in   *kit.SortRequest
		out  *elastic.SortInfo
		err  bool
	}{
		{
			name: "missing first",
			in:   &kit.SortRequest{Field: "field", Asc: false, Missing: kit.SortRequestMissingFirst},
			out:  &elastic.SortInfo{Field: "field", Ascending: false, Missing: EsSortRequestMissingFirst},
			err:  false,
		},
		{
			name: "missing last",
			in:   &kit.SortRequest{Field: "field", Asc: true, Missing: kit.SortRequestMissingLast},
			out:  &elastic.SortInfo{Field: "field", Ascending: true, Missing: EsSortRequestMissingLast},
			err:  false,
		},
		{
			name: "missing empty",
			in:   &kit.SortRequest{Field: "field", Missing: ""},
			out:  &elastic.SortInfo{Field: "field", Ascending: false, Missing: nil},
			err:  false,
		},
		{
			name: "missing not valid",
			in:   &kit.SortRequest{Field: "field", Missing: "notvalidvalue"},
			out:  nil,
			err:  true,
		},
		{
			name: "field empty",
			in:   &kit.SortRequest{Field: ""},
			out:  nil,
			err:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToSortRequestEs(tt.in)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.out, result)
		})
	}
}
