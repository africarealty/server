package impl

import (
	"context"
	"github.com/africarealty/server/src/domain"
	errors "github.com/africarealty/server/src/errors/communications"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/service"
	"github.com/cbroglie/mustache"
)

type templateGeneratorImpl struct {
	storage domain.TemplateStorage
}

func NewTemplateGenerator(storage domain.TemplateStorage) domain.TemplateGenerator {
	t := &templateGeneratorImpl{
		storage: storage,
	}
	return t
}

func (t *templateGeneratorImpl) l() log.CLogger {
	return service.L().Cmp("template-gen")
}

func (t *templateGeneratorImpl) Generate(ctx context.Context, rq *domain.TemplateRequest) (*domain.TemplateResponse, error) {
	t.l().Mth("generate").C(ctx).Dbg()
	if rq.Id == "" {
		return nil, errors.ErrTemplateIdEmpty(ctx)
	}
	// retrieve template
	template, err := t.storage.Get(ctx, rq.Id)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, errors.ErrTemplateNotFound(ctx)
	}
	// render template
	renderedBody, err := mustache.Render(template.Body, rq.Data)
	if err != nil {
		return nil, errors.ErrTemplateGenerator(err, ctx, rq.Id)
	}
	renderedTitle, err := mustache.Render(template.Title, rq.Data)
	if err != nil {
		return nil, errors.ErrTemplateGenerator(err, ctx, rq.Id)
	}
	return &domain.TemplateResponse{
		Title: renderedTitle,
		Body:  renderedBody,
	}, nil
}
