package communications

import (
	"encoding/json"
	aero "github.com/aerospike/aerospike-client-go/v6"
	"github.com/africarealty/server/src/domain"
)

func (t *templateStorageImpl) toTemplateDto(templ *domain.Template) *template {
	if templ == nil {
		return nil
	}
	return &template{
		Id:    templ.Id,
		Title: templ.Title,
		Body:  templ.Body,
	}
}

func (t *templateStorageImpl) toTemplateDomain(templ *template) *domain.Template {
	if templ == nil {
		return nil
	}
	return &domain.Template{
		Id:    templ.Id,
		Title: templ.Title,
		Body:  templ.Body,
	}
}

func (t *templateStorageImpl) toTemplatesDomain(templ []*template) []*domain.Template {
	if templ == nil {
		return nil
	}
	var res []*domain.Template
	for _, tt := range templ {
		res = append(res, t.toTemplateDomain(tt))
	}
	return res
}

func (t *templateStorageImpl) toTemplateCacheDomain(rec *aero.Record) *domain.Template {
	if rec == nil {
		return nil
	}
	body := rec.Bins["template"].(string)
	templ := &domain.Template{}
	_ = json.Unmarshal([]byte(body), templ)
	return templ
}

func (t *templateStorageImpl) toTemplateCache(templ *domain.Template) aero.BinMap {
	templBytes, _ := json.Marshal(templ)
	return aero.BinMap{
		"template": string(templBytes),
	}
}
