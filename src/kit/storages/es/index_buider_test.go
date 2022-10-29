//go:build example
// +build example

package es

import (
	"context"
	"github.com/africarealty/server/src/kit"
	"github.com/africarealty/server/src/kit/log"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
	"testing"
)

func LF() log.CLoggerFunc {
	return func() log.CLogger {
		return log.L(log.Init(&log.Config{Level: log.TraceLevel}))
	}
}

func Test_Index_ChangingModelMappingWithCfgSettings(t *testing.T) {
	es, err := NewEs(&Config{
		Host:     "localhost",
		Port:     "9200",
		Trace:    true,
		Sniff:    true,
		Shards:   2,
		Replicas: 2,
	}, LF())
	if err != nil {
		t.Fatal(err)
	}

	index := kit.NewRandString()
	type model struct {
		Field string `json:"field" es:"type:keyword"`
	}
	err = es.NewBuilder().WithIndex(index).WithMappingModel(&model{}).Build()
	if err != nil {
		t.Fatal(err)
	}

	type modelNew struct {
		Field  string `json:"field" es:"type:keyword"`
		Field2 string `json:"field2" es:"type:keyword"`
	}
	err = es.NewBuilder().WithIndex(index).WithMappingModel(&modelNew{}).Build()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Index_ChangingModelMapping(t *testing.T) {
	es, err := NewEs(&Config{
		Host:  "localhost",
		Port:  "9200",
		Trace: false,
		Sniff: true,
	}, LF())
	if err != nil {
		t.Fatal(err)
	}

	type model struct {
		Field string `json:"field" es:"type:keyword"`
	}
	index := kit.NewRandString()
	err = es.NewBuilder().WithIndex(index).WithMappingModel(&model{}).Build()
	if err != nil {
		t.Fatal(err)
	}
	type modelNew struct {
		Field  string `json:"field" es:"type:keyword"`
		Field2 string `json:"field2" es:"type:keyword"`
	}

	err = es.NewBuilder().WithIndex(index).WithMappingModel(&modelNew{}).Build()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Index_ChangingMapping_ExistentFields(t *testing.T) {
	es, err := NewEs(&Config{
		Host:  "localhost",
		Port:  "9200",
		Trace: false,
		Sniff: true,
	}, LF())
	if err != nil {
		t.Fatal(err)
	}

	type model struct {
		Field string `json:"field" es:"type:keyword"`
	}

	index := kit.NewRandString()
	err = es.NewBuilder().WithIndex(index).WithMappingModel(&model{}).Build()
	if err != nil {
		t.Fatal(err)
	}
	type modelNew struct {
		Field string `json:"field" es:"type:text"`
	}
	err = es.NewBuilder().WithIndex(index).WithMappingModel(&modelNew{}).Build()
	assert.NotNil(t, err)
}

func Test_Index_ChangingMapping_ExplicitMapping(t *testing.T) {
	es, err := NewEs(&Config{
		Host:     "localhost",
		Port:     "9200",
		Trace:    true,
		Sniff:    true,
		Shards:   2,
		Replicas: 2,
	}, LF())
	if err != nil {
		t.Fatal(err)
	}

	mapping := `
{
	"mappings": {
		"properties": {
			"field1": {
				"type": "keyword"
			}
		}
	}
}
`
	index := kit.NewRandString()
	err = es.NewBuilder().WithIndex(index).WithExplicitMapping(mapping).Build()
	if err != nil {
		t.Fatal(err)
	}
	newMapping := `
{
	"mappings": {
		"properties": {
			"field1": {
				"type": "keyword"
			},
			"field2": {
				"type": "keyword"
			}
		}
	}
}
`
	err = es.NewBuilder().WithIndex(index).WithExplicitMapping(newMapping).Build()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Index_Mapping_WhenNotIndexField(t *testing.T) {
	es, err := NewEs(&Config{
		Host:  "localhost",
		Port:  "9200",
		Trace: false,
		Sniff: true,
	}, LF())
	if err != nil {
		t.Fatal(err)
	}

	type model struct {
		Field1 string `json:"field1" es:"type:keyword"`
		Field2 string `json:"field2" es:"type:keyword;-"`
	}

	index := kit.NewRandString()
	err = es.NewBuilder().WithIndex(index).WithMappingModel(&model{}).Build()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Alias_ChangingModelMapping(t *testing.T) {
	es, err := NewEs(&Config{
		Host:     "localhost",
		Port:     "9200",
		Trace:    true,
		Sniff:    true,
		Shards:   2,
		Replicas: 2,
	}, LF())
	if err != nil {
		t.Fatal(err)
	}

	// create alias and index
	alias := kit.NewRandString()
	type model struct {
		Field string `json:"field" es:"type:keyword"`
	}
	err = es.NewBuilder().WithAlias(alias).WithMappingModel(&model{}).Build()
	if err != nil {
		t.Fatal(err)
	}

	// index document
	err = es.Index(alias, kit.NewId(), &model{Field: "value"})
	if err != nil {
		t.Fatal(err)
	}

	// search
	err = es.Refresh(alias)
	if err != nil {
		t.Fatal(err)
	}

	// get data from alias
	srchRs, err := es.GetClient().Search().Index(alias).Query(elastic.NewMatchAllQuery()).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, int64(1), srchRs.TotalHits())

	// change mapping
	type modelNew struct {
		Field  string `json:"field" es:"type:keyword"`
		Field2 string `json:"field2" es:"type:keyword"`
	}
	err = es.NewBuilder().WithAlias(alias).WithMappingModel(&modelNew{}).Build()
	if err != nil {
		t.Fatal(err)
	}

	// index document
	err = es.Index(alias, kit.NewId(), &modelNew{Field: "value", Field2: "value2"})
	if err != nil {
		t.Fatal(err)
	}

	// search
	err = es.Refresh(alias)
	if err != nil {
		t.Fatal(err)
	}

	// get data from alias
	srchRs, err = es.GetClient().Search().Index(alias).Query(elastic.NewMatchAllQuery()).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, int64(2), srchRs.TotalHits())
}
