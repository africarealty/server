package communications

import (
	"context"
	"github.com/africarealty/server/src/domain"
	errors "github.com/africarealty/server/src/errors/communications"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/service"
)

type templateImpl struct {
	storage domain.TemplateStorage
}

func NewTemplateService(storage domain.TemplateStorage) domain.TemplateService {
	return &templateImpl{
		storage: storage,
	}
}

func (t *templateImpl) l() log.CLogger {
	return service.L().Cmp("template")
}

func (t *templateImpl) CreateTemplate(ctx context.Context, rq *domain.Template) (*domain.Template, error) {
	t.l().Mth("create").C(ctx).Dbg()
	if rq.Id == "" {
		return nil, errors.ErrTemplateIdEmpty(ctx)
	}
	if rq.Title == "" {
		return nil, errors.ErrTemplateTitleEmpty(ctx)
	}
	if rq.Body == "" {
		return nil, errors.ErrTemplateBodyEmpty(ctx)
	}

	// get stored template
	templateStored, err := t.storage.Get(ctx, rq.Id)
	if err != nil {
		return nil, err
	}
	if templateStored != nil {
		return nil, errors.ErrTemplateAlreadyExists(ctx)
	}

	template := &domain.Template{
		Id:    rq.Id,
		Title: rq.Title,
		Body:  rq.Body,
	}

	err = t.storage.Create(ctx, template)
	if err != nil {
		return nil, err
	}
	return template, nil
}

func (t *templateImpl) UpdateTemplate(ctx context.Context, template *domain.Template) (*domain.Template, error) {
	t.l().Mth("update").C(ctx).Dbg()

	if template.Id == "" {
		return nil, errors.ErrTemplateIdEmpty(ctx)
	}
	if template.Title == "" {
		return nil, errors.ErrTemplateTitleEmpty(ctx)
	}
	if template.Body == "" {
		return nil, errors.ErrTemplateBodyEmpty(ctx)
	}

	// get stored template
	templateStored, err := t.storage.Get(ctx, template.Id)
	if err != nil {
		return nil, err
	}
	if templateStored == nil {
		return nil, errors.ErrTemplateNotFound(ctx)
	}

	// save to store
	err = t.storage.Update(ctx, template)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (t *templateImpl) DeleteTemplate(ctx context.Context, id string) error {
	t.l().Mth("delete").C(ctx).Dbg()

	if id == "" {
		return errors.ErrTemplateIdEmpty(ctx)
	}

	// get stored template
	template, err := t.storage.Get(ctx, id)
	if err != nil {
		return err
	}
	if template == nil {
		return errors.ErrTemplateNotFound(ctx)
	}

	return t.storage.Delete(ctx, id)
}

func (t *templateImpl) GetTemplate(ctx context.Context, id string) (*domain.Template, error) {
	t.l().Mth("get").C(ctx).Dbg()

	if id == "" {
		return nil, errors.ErrTemplateIdEmpty(ctx)
	}

	// get stored template
	template, err := t.storage.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, errors.ErrTemplateNotFound(ctx)
	}

	return template, nil
}

func (t *templateImpl) SearchTemplates(ctx context.Context, query string) ([]*domain.Template, error) {
	t.l().Mth("search").C(ctx).Dbg()

	// search templates
	templates, err := t.storage.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	if templates == nil {
		return nil, errors.ErrTemplateNotFound(ctx)
	}

	return templates, nil
}
