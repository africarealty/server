package storage

import (
	"context"
	aero "github.com/aerospike/aerospike-client-go/v6"
	"github.com/aerospike/aerospike-client-go/v6/types"
	"github.com/africarealty/server/src/errors/auth"
	"github.com/africarealty/server/src/kit/auth"
	"github.com/africarealty/server/src/kit/goroutine"
	"github.com/africarealty/server/src/kit/log"
	kitAero "github.com/africarealty/server/src/kit/storages/aerospike"
	"github.com/africarealty/server/src/kit/storages/pg"
	"github.com/africarealty/server/src/service"
	"time"
)

type sessionDetails struct {
	Roles []string `json:"roles,omitempty"`
}

type session struct {
	Id             string     `gorm:"column:id"`
	UserId         string     `gorm:"column:user_id"`
	Username       string     `gorm:"column:username"`
	LoginAt        time.Time  `gorm:"column:login_at"`
	LastActivityAt time.Time  `gorm:"column:last_activity_at"`
	LogoutAt       *time.Time `gorm:"column:logout_at"`
	Details        string     `gorm:"column:details"`
}

type SessionStorageImpl struct {
	pg   *pg.Storage
	aero kitAero.Aerospike
	cfg  *kitAero.Config
}

func (s *SessionStorageImpl) l() log.CLogger {
	return service.L().Cmp("session-storage")
}

func NewSessionStorage(pg *pg.Storage, aero kitAero.Aerospike, cfg *kitAero.Config) *SessionStorageImpl {
	return &SessionStorageImpl{
		pg:   pg,
		aero: aero,
		cfg:  cfg,
	}
}

func (s *SessionStorageImpl) getFromCacheById(ctx context.Context, sid string) (*auth.Session, error) {
	s.l().Mth("get-cache").C(ctx).F(log.FF{"sid": sid}).Trc()
	key, err := aero.NewKey(AeroNsCache, AeroSetSessionCache, sid)
	if err != nil {
		return nil, errors.ErrSessionStorageAeroKey(err, ctx)
	}
	policy := aero.NewPolicy()
	policy.SendKey = true
	rec, err := s.aero.Instance().Get(policy, key)
	if err != nil && !err.Matches(types.KEY_NOT_FOUND_ERROR) {
		return nil, errors.ErrSessionStorageGetCache(err, ctx)
	}
	return s.toSessionCacheDomain(rec), nil
}

func (s *SessionStorageImpl) setCache(ctx context.Context, session *auth.Session) error {
	s.l().Mth("set-cache").C(ctx).F(log.FF{"sid": session.Id}).Trc()
	key, err := aero.NewKey(AeroNsCache, AeroSetSessionCache, session.Id)
	if err != nil {
		return errors.ErrSessionStorageAeroKey(err, ctx)
	}
	writePolicy := aero.NewWritePolicy(0, 1800)
	writePolicy.SendKey = true
	err = s.aero.Instance().Put(writePolicy, key, s.toSessionCache(session))
	if err != nil {
		return errors.ErrSessionStoragePutCache(err, ctx)
	}
	return nil
}

func (s *SessionStorageImpl) clearCache(ctx context.Context, sid string) error {
	s.l().Mth("clear-cache").C(ctx).F(log.FF{"sid": sid}).Trc()
	key, err := aero.NewKey(AeroNsCache, AeroSetSessionCache, sid)
	if err != nil {
		return errors.ErrSessionStorageAeroKey(err, ctx)
	}
	_, err = s.aero.Instance().Delete(nil, key)
	if err != nil {
		return errors.ErrSessionStorageClearCache(err, ctx)
	}
	return nil
}

func (s *SessionStorageImpl) Get(ctx context.Context, sid string) (*auth.Session, error) {
	l := s.l().Mth("get").C(ctx).F(log.FF{"sid": sid}).Trc()
	if sid == "" {
		return nil, nil
	}
	// check cache first
	sess, err := s.getFromCacheById(ctx, sid)
	if err != nil {
		return nil, err
	}
	if sess != nil {
		l.Trc("found in cache")
		return sess, nil
	}
	// get from db
	dto := &session{Id: sid}
	res := s.pg.Instance.Limit(1).Find(&dto)
	if res.Error != nil {
		return nil, errors.ErrSessionStorageGetDb(res.Error, ctx)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	sess = s.toSessionDomain(dto)
	// set cache
	err = s.setCache(ctx, sess)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *SessionStorageImpl) GetByUser(ctx context.Context, uid string) ([]*auth.Session, error) {
	s.l().C(ctx).Mth("get-by-user").F(log.FF{"uid": uid}).Trc()
	if uid == "" {
		return []*auth.Session{}, nil
	}
	var sessions []*session
	if err := s.pg.Instance.Where("user_id = ?::uuid and logout_at is null", uid).Find(&sessions).Error; err == nil {
		return s.toSessionsDomain(sessions), nil
	} else {
		return nil, errors.ErrSessionGetByUser(err, ctx)
	}
}

func (s *SessionStorageImpl) CreateSession(ctx context.Context, session *auth.Session) error {
	l := s.l().C(ctx).Mth("create").F(log.FF{"sid": session.Id}).Trc()
	eg := goroutine.NewGroup(ctx).WithLogger(l)
	// session
	eg.Go(func() error {
		if err := s.pg.Instance.Create(s.toSessionDto(session)).Error; err != nil {
			return errors.ErrSessionStorageCreateSession(err, ctx)
		}
		l.Dbg("session created")
		return nil
	})
	// set cache
	eg.Go(func() error {
		if err := s.setCache(ctx, session); err != nil {
			return err
		}
		return nil
	})
	return eg.Wait()
}

func (s *SessionStorageImpl) UpdateLastActivity(ctx context.Context, sid string, lastActivity time.Time) error {
	s.l().Mth("logout").C(ctx).F(log.FF{"sid": sid}).Dbg()
	// update DB
	if err := s.pg.Instance.Model(&session{Id: sid}).
		Updates(map[string]interface{}{
			"last_activity_at": lastActivity,
		}).Error; err != nil {
		return errors.ErrSessionStorageUpdateLastActivity(err, ctx)
	}
	return nil
}

func (s *SessionStorageImpl) Logout(ctx context.Context, sid string, logoutAt time.Time) error {
	s.l().Mth("logout").C(ctx).F(log.FF{"sid": sid}).Trc()
	if err := s.pg.Instance.Model(&session{Id: sid}).
		Updates(map[string]interface{}{
			"logout_at": logoutAt,
		}).Error; err != nil {
		return errors.ErrSessionStorageUpdateLogout(err, ctx)
	}
	// clear cache
	if err := s.clearCache(ctx, sid); err != nil {
		return err
	}
	return nil
}
