package es

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/africarealty/server/src/kit"
	"github.com/africarealty/server/src/kit/log"
	"github.com/olivere/elastic/v7"
	"reflect"
	"strings"
)

// IndexBuilder allows creating / modification a ES index
type IndexBuilder interface {
	// WithAlias specifies an index with alias
	WithAlias(name string) IndexBuilder
	// WithIndex specifies an index name
	// if alias specified with WithAlias call, you don't need specify index name explicitly
	WithIndex(name string) IndexBuilder
	// WithMappingModel specifies index mapping based on model provided
	// if index doesn't exist, a new index is created
	// if index exists it checks whether existent mapping is modified and if it is, it fails. If only new fields added, it handles them as PUT
	// Note, "json" tag must be specified together with "es" tag
	// example:
	// type IndexModel struct {
	//   Field1 string `json:"field1" es:"type:text"`   	// field is mapped with text type
	//   Field2 string `json:"field2" es:"type:keyword"` 	// field is mapped with keyword type
	//   Field3 time.Time `json:"field3" es:"type:date"` 	// field is mapped with date type
	//   Field4 time.Time `json:"field4" es:"-"` 			// field is mapped with "Index=false"
	// }
	//
	// model must be pointer type
	WithMappingModel(model interface{}) IndexBuilder
	// WithMappingModel specifies index mapping explicitly as a serialized json mapping object
	// see ES doc about https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping.html
	// it checks if index exists, if not creates it
	WithExplicitMapping(mapping string) IndexBuilder
	// Build builds a new alias/index or modifies mapping of an existent index
	Build() error
}

type esIndexBuilder struct {
	client          *elastic.Client
	logger          log.CLoggerFunc
	cfg             *Config
	alias           string
	index           string
	mappingModel    interface{}
	mappingExplicit string
}

func (e *esIndexBuilder) l() log.CLogger {
	return e.logger().Cmp("es-idx-builder")
}

func (e *esIndexBuilder) WithAlias(name string) IndexBuilder {
	e.alias = name
	return e
}

func (e *esIndexBuilder) WithIndex(name string) IndexBuilder {
	e.index = name
	return e
}

func (e *esIndexBuilder) WithMappingModel(model interface{}) IndexBuilder {
	e.mappingModel = model
	return e
}

func (e *esIndexBuilder) WithExplicitMapping(mapping string) IndexBuilder {
	e.mappingExplicit = mapping
	return e
}

func (e *esIndexBuilder) getMapping(index string) (*EsMapping, error) {
	// get current mapping
	curMappings, err := e.client.GetMapping().Index(index).Do(context.Background())
	if err != nil {
		return nil, ErrEsGetMapping(err, index)
	}
	curMapping, ok := curMappings[index]
	if !ok {
		return nil, ErrEsNoMappingFound(index)
	}

	mappingJson, _ := json.Marshal(curMapping)
	currentMapping := &EsMapping{}
	err = json.Unmarshal(mappingJson, currentMapping)
	if err != nil {
		return nil, ErrEsMappingSchemaNotExpected(err, index)
	}

	return currentMapping, nil
}

// modelToMapping creates ES mapping based on model tag
// check model_mapping_test for usage details
func (e *esIndexBuilder) modelToMapping(modelObj interface{}) (*EsMapping, error) {

	e.l().Mth("model-to-mapping").Dbg()

	if modelObj == nil {
		return nil, nil
	}

	type params map[string]string

	if reflect.ValueOf(modelObj).Kind() != reflect.Ptr || reflect.TypeOf(modelObj).Elem().Kind() != reflect.Struct {
		return nil, ErrEsInvalidModelType()
	}

	// takes type description
	r := reflect.TypeOf(modelObj).Elem()
	mappingProperties := make(EsProperties)

	// build mapping fields map
	// go through fields
	for i := 0; i < r.NumField(); i++ {
		field := r.Field(i)

		// check json tag
		// we use index field name from json mapping
		// if json tag missing, field is skipped
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {

			jsonParams := strings.Split(jsonTag, ",")
			// if there is no field name in json tag
			if len(jsonParams) == 0 {
				return nil, ErrEsInvalidModel()
			}
			indexFieldName := jsonParams[0]

			// take es tag
			esTag := field.Tag.Get("es")
			// if es tag missing, skip the field
			if esTag != "" {
				esTagParams := make(params)
				// take params separated by ;
				params := strings.Split(esTag, ";")
				for _, p := range params {
					kv := strings.Split(p, ":")
					if len(kv) == 2 {
						esTagParams[kv[0]] = kv[1]
					} else {
						esTagParams[kv[0]] = ""
					}
				}
				// populate mapping params
				mappingProperties[indexFieldName] = &EsProperty{}

				if _, ok := esTagParams["-"]; ok {
					// if empty sign exists
					f := false
					mappingProperties[indexFieldName].Type = EsTypeKeyword
					mappingProperties[indexFieldName].Index = &f
				} else if esTypeTag, ok := esTagParams["type"]; ok {
					if _, ok := typesMap[esTypeTag]; !ok {
						return nil, ErrEsInvalidModel()
					} else {
						mappingProperties[indexFieldName].Type = esTypeTag
					}
				}

			}
		}
	}

	// return ES mapping if specified
	if len(mappingProperties) == 0 {
		return nil, ErrEsInvalidModel()
	} else {
		r := &EsMapping{}
		r.Mappings.Properties = mappingProperties
		return r, nil
	}

}

func (e *esIndexBuilder) getNewMapping(index string) (*EsMapping, error) {
	// get new mapping
	newMapping := &EsMapping{}
	var err error
	if e.mappingModel != nil {
		// if mapping specified as a model
		newMapping, err = e.modelToMapping(e.mappingModel)
		if err != nil {
			return nil, err
		}
	} else {
		// if mapping specified explicitly
		err = json.Unmarshal([]byte(e.mappingExplicit), newMapping)
		if err != nil {
			return nil, ErrEsMappingSchemaNotExpected(err, index)
		}
	}
	return newMapping, nil
}

func (e *esIndexBuilder) modifyMapping(index string, curMapping, newMapping *EsMapping) error {
	l := e.l().Mth("modify-mapping").Dbg()
	// check if there are changes in existent fields
	if v := e.checkExistentFieldsMappingModified(curMapping, newMapping); len(v) > 0 {
		return ErrEsMappingExistentFieldsModified(index, v)
	}
	// extract added fields
	if addedFieldsMapping := e.addedFieldsMapping(curMapping, newMapping); len(addedFieldsMapping.Mappings.Properties) > 0 {
		addedFieldsMappingJson, _ := json.Marshal(addedFieldsMapping.Mappings)
		_, err := e.client.PutMapping().Index(index).BodyString(string(addedFieldsMappingJson)).Do(context.Background())
		if err != nil {
			return ErrEsPutMapping(err, index)
		}
		l.DbgF("fields added: %+v", addedFieldsMapping.Mappings.Properties)
	}
	return nil
}

func (e *esIndexBuilder) createIndex(index string, mapping *EsMapping) error {
	l := e.l().Mth("create-index").F(log.FF{"index": index}).Dbg()
	// set index settings
	e.setSettings(mapping)
	// prepare mapping
	newMappingJson, _ := json.Marshal(mapping)
	// create
	_, err := e.client.CreateIndex(index).BodyString(string(newMappingJson)).Do(context.Background())
	if err != nil {
		return ErrEsIdxCreate(err, index)
	}
	l.Dbg("created")
	return nil
}

func (e *esIndexBuilder) buildAlias(alias string) error {
	l := e.l().Mth("build-alias").F(log.FF{"alias": alias}).Dbg()
	// check alias exists
	exists, err := e.client.IndexExists(alias).Do(context.Background())
	if err != nil {
		return ErrEsIdxExists(err, alias)
	}
	if exists {
		// existent alias

		// we allow adding new fields to mapping, but don't allow changing existent ones
		l.DbgF("alias %s exists", alias)

		// get indexes by alias
		aliasesRs, err := e.client.Aliases().Alias(alias).Do(context.Background())
		if err != nil {
			return ErrEsGetIndexesByAlias(err, alias)
		}
		if len(aliasesRs.Indices) == 0 {
			return ErrEsNoIndicesForAlias(alias)
		}
		// get writable index
		var writeIndexName string
	loop:
		for idxName, idx := range aliasesRs.Indices {
			for _, ia := range idx.Aliases {
				if ia.IsWriteIndex {
					writeIndexName = idxName
					break loop
				}
			}
		}
		if writeIndexName == "" {
			return ErrEsNoWriteIndexForAlias(alias)
		}
		l.F(log.FF{"writeIndex": writeIndexName})
		// get current mapping
		currentMapping, err := e.getMapping(writeIndexName)
		if err != nil {
			return err
		}
		// get new mapping
		newMapping, err := e.getNewMapping(writeIndexName)
		if err != nil {
			return err
		}
		// modify mapping for alias (it modifies mapping for all the indexes)
		err = e.modifyMapping(alias, currentMapping, newMapping)
		if err != nil {
			return err
		}
		l.Dbg("modified")
	} else {
		// new alias

		// get new mapping
		newMapping, err := e.getNewMapping(alias)
		if err != nil {
			return err
		}
		// create write index
		idxName := fmt.Sprintf("%s-idx-%s", alias, kit.Now().Format("20060102150405"))
		err = e.createIndex(idxName, newMapping)
		if err != nil {
			return err
		}
		// add index to alias
		_, err = e.client.Alias().
			Action(elastic.NewAliasAddAction(alias).Index(idxName).IsWriteIndex(true)).
			Do(context.Background())
		if err != nil {
			return err
		}
		l.Dbg("created")
	}
	return nil
}

func (e *esIndexBuilder) buildIndex(index string) error {
	l := e.l().Mth("build-index").F(log.FF{"index": index}).Dbg()
	// check index exists
	exists, err := e.client.IndexExists(index).Do(context.Background())
	if err != nil {
		return ErrEsIdxExists(err, index)
	}
	if exists {
		// get current mapping
		currentMapping, err := e.getMapping(index)
		if err != nil {
			return err
		}
		// get new mapping
		newMapping, err := e.getNewMapping(index)
		if err != nil {
			return err
		}
		// modify mapping for index
		err = e.modifyMapping(index, currentMapping, newMapping)
		if err != nil {
			return err
		}
		l.Dbg("modified")
	} else {
		// new index

		// get new mapping
		newMapping, err := e.getNewMapping(index)
		if err != nil {
			return err
		}
		// create write index
		err = e.createIndex(index, newMapping)
		if err != nil {
			return err
		}
		l.Dbg("created")
	}
	return nil
}

func (e *esIndexBuilder) Build() error {
	e.l().Mth("build").Dbg()

	// check passed params
	if e.alias == "" && e.index == "" {
		return ErrEsIndexBuilderAliasAndIndexEmpty()
	}
	if e.mappingExplicit == "" && e.mappingModel == nil {
		return ErrEsIndexBuilderModelEmpty()
	}

	// alias-based
	if e.alias != "" {
		return e.buildAlias(e.alias)
	} else {
		return e.buildIndex(e.index)
	}

}

func (e *esIndexBuilder) setSettings(mapping *EsMapping) {
	if mapping.Settings.NumberOfReplicas == 0 {
		mapping.Settings.NumberOfReplicas = e.cfg.Replicas
		if mapping.Settings.NumberOfReplicas == 0 {
			mapping.Settings.NumberOfReplicas = 1
		}
	}
	if mapping.Settings.NumberOfShards == 0 {
		mapping.Settings.NumberOfShards = e.cfg.Shards
		if mapping.Settings.NumberOfShards == 0 {
			mapping.Settings.NumberOfShards = 1
		}
	}
}

func (e *esIndexBuilder) addedFieldsMapping(currentMapping, newMapping *EsMapping) *EsMapping {
	addedFieldsMapping := &EsMapping{}
	addedFieldsMapping.Mappings.Properties = make(EsProperties)
	for f, v := range newMapping.Mappings.Properties {
		if _, found := currentMapping.Mappings.Properties[f]; !found {
			addedFieldsMapping.Mappings.Properties[f] = v
		}
	}
	return addedFieldsMapping
}

// checkExistentFieldsMappingModified compares current and provided mapping and returns true if there are changes in existent fields
func (e *esIndexBuilder) checkExistentFieldsMappingModified(currentMapping, newMapping *EsMapping) []string {
	var modifiedFields []string
	for curFieldName, curField := range currentMapping.Mappings.Properties {
		for newFieldName, newField := range newMapping.Mappings.Properties {
			if curFieldName == newFieldName && curField.Type != newField.Type {
				modifiedFields = append(modifiedFields, curFieldName)
			}
		}
	}
	return modifiedFields
}
