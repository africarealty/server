package communications

import (
	"context"
	"github.com/africarealty/server/src/domain"
	errors "github.com/africarealty/server/src/errors/communications"
	"github.com/africarealty/server/src/kit"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/kit/storages/pg"
	"github.com/africarealty/server/src/service"
)

type email struct {
	pg.GormDto
	Id               string  `gorm:"column:id"`
	UserId           *string `gorm:"column:user_id"`
	Email            string  `gorm:"column:email"`
	Text             string  `gorm:"column:text"`
	TemplateId       *string `gorm:"column:template_id"`
	TemplateData     *string `gorm:"column:template_data"`
	SendStatus       string  `gorm:"column:send_status"`
	ErrorDescription *string `gorm:"column:error_desc"`
	Attachments      *string `gorm:"attachments"`
}

type EmailStorageImpl struct {
	pg *pg.Storage
}

func NewEmailStorage(pg *pg.Storage) *EmailStorageImpl {
	return &EmailStorageImpl{
		pg: pg,
	}
}

func (e *EmailStorageImpl) l() log.CLogger {
	return service.L().Cmp("email-storage")
}

func (e *EmailStorageImpl) UpdateEmail(ctx context.Context, request *domain.Email) error {
	e.l().C(ctx).Mth("update").Dbg()
	t := kit.Now()
	dto := e.toEmailDto(request)
	dto.UpdatedAt = &t
	result := e.pg.Instance.Exec(`
			update emails set
				send_status = ?,
				error_desc = ?,
				updated_at = ?
			where id = ?::uuid 
			`, dto.SendStatus, dto.ErrorDescription, dto.UpdatedAt, dto.Id)
	if result.Error != nil {
		return errors.ErrEmailStorageUpdateEmailDb(result.Error, ctx)
	}
	return nil
}

func (e *EmailStorageImpl) CreateEmail(ctx context.Context, request *domain.Email) error {
	e.l().C(ctx).Mth("create").Dbg()
	dto := e.toEmailDto(request)
	result := e.pg.Instance.Create(&dto)
	if result.Error != nil {
		return errors.ErrEmailStorageCreateEmailDb(result.Error, ctx)
	}
	return nil
}
