package es

import "github.com/africarealty/server/src/kit/er"

var (
	ErrCodeEsNewClient                      = "ES-001"
	ErrCodeEsIdxExists                      = "ES-002"
	ErrCodeEsIdx                            = "ES-003"
	ErrCodeEsIdxAsync                       = "ES-004"
	ErrCodeEsIdxCreate                      = "ES-007"
	ErrCodeEsBulkIdx                        = "ES-008"
	ErrCodeEsExists                         = "ES-009"
	ErrCodeEsInvalidModel                   = "ES-011"
	ErrCodeEsInvalidModelType               = "ES-012"
	ErrCodeEsGetMapping                     = "ES-013"
	ErrCodeEsNoMappingFound                 = "ES-014"
	ErrCodeEsMappingSchemaNotExpected       = "ES-015"
	ErrCodeEsMappingExistentFieldsModified  = "ES-016"
	ErrCodeEsPutMapping                     = "ES-017"
	ErrCodeEsDel                            = "ES-018"
	ErrCodeEsIndexBuilderAliasAndIndexEmpty = "ES-019"
	ErrCodeEsIndexBuilderModelEmpty         = "ES-020"
	ErrCodeEsGetIndexesByAlias              = "ES-021"
	ErrCodeEsNoIndicesForAlias              = "ES-022"
	ErrCodeEsNoWriteIndexForAlias           = "ES-023"
	ErrCodeEsRefresh                        = "ES-024"
	ErrCodeEsBasicAuthInvalid               = "ES-025"
	ErrCodeEsBulkDel                        = "ES-026"
	ErrCodeEsSortRequestMissingInvalid      = "ES-027"
	ErrCodeEsSortRequestFieldEmpty          = "ES-028"
)

var (
	ErrEsNewClient = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeEsNewClient, "").Err() }
	ErrEsIdxExists = func(cause error, index string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsIdxExists, "").F(er.FF{"idx": index}).Err()
	}
	ErrEsGetMapping = func(cause error, index string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsGetMapping, "").F(er.FF{"idx": index}).Err()
	}
	ErrEsGetIndexesByAlias = func(cause error, alias string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsGetIndexesByAlias, "").F(er.FF{"alias": alias}).Err()
	}
	ErrEsMappingSchemaNotExpected = func(cause error, index string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsMappingSchemaNotExpected, "").F(er.FF{"idx": index}).Err()
	}
	ErrEsPutMapping = func(cause error, index string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsPutMapping, "").F(er.FF{"idx": index}).Err()
	}
	ErrEsIdx = func(cause error, index, id string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsIdx, "").F(er.FF{"idx": index, "id": id}).Err()
	}
	ErrEsDel = func(cause error, index, id string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsDel, "").F(er.FF{"idx": index, "id": id}).Err()
	}
	ErrEsRefresh = func(cause error, index string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsRefresh, "").F(er.FF{"idx": index}).Err()
	}
	ErrEsBulkIdx = func(cause error, index string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsBulkIdx, "").F(er.FF{"idx": index}).Err()
	}
	ErrEsBulkDel = func(cause error, index string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsBulkDel, "").F(er.FF{"idx": index}).Err()
	}
	ErrEsIdxCreate = func(cause error, index string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsIdxCreate, "").F(er.FF{"idx": index}).Err()
	}
	ErrEsExists = func(cause error, index, id string) error {
		return er.WrapWithBuilder(cause, ErrCodeEsExists, "").F(er.FF{"idx": index, "id": id}).Err()
	}
	ErrEsNoMappingFound = func(index string) error {
		return er.WithBuilder(ErrCodeEsNoMappingFound, "no mapping found").F(er.FF{"idx": index}).Err()
	}
	ErrEsMappingExistentFieldsModified = func(index string, fields []string) error {
		return er.WithBuilder(ErrCodeEsMappingExistentFieldsModified, "ES doesn't allow changing mapping for existent fields.").F(er.FF{"idx": index, "fields": fields}).Err()
	}
	ErrEsInvalidModel     = func() error { return er.WithBuilder(ErrCodeEsInvalidModel, "invalid model, check tags").Err() }
	ErrEsInvalidModelType = func() error {
		return er.WithBuilder(ErrCodeEsInvalidModelType, "model must be pointer of struct").Err()
	}
	ErrEsIndexBuilderAliasAndIndexEmpty = func() error {
		return er.WithBuilder(ErrCodeEsIndexBuilderAliasAndIndexEmpty, "neither alias name nor index name specified").Err()
	}
	ErrEsIndexBuilderModelEmpty = func() error {
		return er.WithBuilder(ErrCodeEsIndexBuilderModelEmpty, "model not specified").Err()
	}
	ErrEsNoIndicesForAlias = func(alias string) error {
		return er.WithBuilder(ErrCodeEsNoIndicesForAlias, "model not specified").F(er.FF{"alias": alias}).Err()
	}
	ErrEsNoWriteIndexForAlias = func(alias string) error {
		return er.WithBuilder(ErrCodeEsNoWriteIndexForAlias, "no write index").F(er.FF{"alias": alias}).Err()
	}
	ErrEsBasicAuthInvalid = func() error {
		return er.WithBuilder(ErrCodeEsBasicAuthInvalid, "basic auth invalid").Err()
	}
	ErrEsSortRequestMissingInvalid = func(missing string) error {
		return er.WithBuilder(ErrCodeEsSortRequestMissingInvalid, "sort request missing parameter invalid").F(er.FF{"missing": missing}).Err()
	}
	ErrEsSortRequestFieldEmpty = func() error {
		return er.WithBuilder(ErrCodeEsSortRequestFieldEmpty, "sort request field parameter empty").Err()
	}
)
