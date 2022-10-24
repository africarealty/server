package communications

import (
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/kit"
	"github.com/africarealty/server/src/kit/storages/pg"
)

func (e *EmailStorageImpl) toEmailDto(d *domain.Email) *email {
	if d == nil {
		return nil
	}

	var templData *string
	var templId string
	if d.Template != nil {
		if d.Template.Data != nil {
			v := kit.Json(d.Template.Data)
			templData = &v
		}
		templId = d.Template.Id
	}

	var attachments *string
	if len(d.LinkFileIds) > 0 {
		v := kit.Json(d.LinkFileIds)
		attachments = &v
	}

	return &email{
		Id:               d.Id,
		UserId:           pg.StringToNull(d.UserId),
		Email:            d.Email,
		Text:             d.Text,
		TemplateId:       pg.StringToNull(templId),
		TemplateData:     templData,
		SendStatus:       d.SendStatus,
		ErrorDescription: pg.StringToNull(d.ErrorDescription),
		Attachments:      attachments,
	}
}
