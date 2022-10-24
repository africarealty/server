package auth

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
	"github.com/africarealty/server/src/repository/storage"
	"github.com/africarealty/server/src/service"
	"time"
)

const (
	AeroSetUserCache = "user_cache"
)

type userDetails struct {
	FirstName string   `json:"firstName,omitempty"`
	LastName  string   `json:"lastName,omitempty"`
	Groups    []string `json:"groups,omitempty"`
	Roles     []string `json:"roles,omitempty"`
}

type user struct {
	pg.GormDto
	Id          string     `gorm:"column:id"`
	Username    string     `gorm:"column:username"`
	Password    *string    `gorm:"column:password"`
	Type        string     `gorm:"column:type"`
	Details     string     `gorm:"column:details"`
	ActivatedAt *time.Time `gorm:"column:activated_at"`
	LockedAt    *time.Time `gorm:"column:locked_at"`
}

type UserStorageImpl struct {
	pg   *pg.Storage
	aero kitAero.Aerospike
	cfg  *kitAero.Config
}

func (s *UserStorageImpl) l() log.CLogger {
	return service.L().Cmp("user-storage")
}

func NewUserStorage(pg *pg.Storage, aero kitAero.Aerospike, cfg *kitAero.Config) *UserStorageImpl {
	return &UserStorageImpl{
		pg:   pg,
		aero: aero,
		cfg:  cfg,
	}
}

func (s *UserStorageImpl) init(ctx context.Context) error {
	s.l().Mth("init").C(ctx).Trc()
	// create secondary index
	task, err := s.aero.Instance().CreateIndex(nil, storage.AeroNsCache, AeroSetUserCache, "idx_un", "username", aero.STRING)
	if err != nil && !err.Matches(types.INDEX_FOUND) {
		return errors.ErrUserStorageCreateIndex(err, ctx)
	}
	if task != nil {
		if err := <-task.OnComplete(); err != nil {
			return errors.ErrUserStorageCreateIndex(err, ctx)
		}
	}
	return nil
}

func (s *UserStorageImpl) clearCache(ctx context.Context, userId string) error {
	s.l().Mth("clear-cache").C(ctx).F(log.FF{"userId": userId}).Trc()
	key, err := aero.NewKey(storage.AeroNsCache, AeroSetUserCache, userId)
	if err != nil {
		return errors.ErrUserStorageAeroKey(err, ctx)
	}
	_, err = s.aero.Instance().Delete(nil, key)
	if err != nil {
		return errors.ErrUserStorageClearCache(err, ctx)
	}
	return nil
}

func (s *UserStorageImpl) getFromCacheById(ctx context.Context, userId string) (*auth.User, error) {
	s.l().Mth("get-cache").C(ctx).F(log.FF{"userId": userId}).Trc()
	key, err := aero.NewKey(storage.AeroNsCache, AeroSetUserCache, userId)
	if err != nil {
		return nil, errors.ErrUserStorageAeroKey(err, ctx)
	}
	policy := aero.NewPolicy()
	policy.SendKey = true
	rec, err := s.aero.Instance().Get(policy, key)
	if err != nil && !err.Matches(types.KEY_NOT_FOUND_ERROR) {
		return nil, errors.ErrUserStorageGetCache(err, ctx)
	}
	return s.toUserCacheDomain(rec), nil
}

func (s *UserStorageImpl) getFromCacheByUsername(ctx context.Context, username string) (*auth.User, error) {
	s.l().Mth("get-cache").C(ctx).F(log.FF{"username": username}).Trc()
	queryPolicy := aero.NewQueryPolicy()
	queryPolicy.SendKey = true
	queryPolicy.MaxRecords = int64(1)
	queryPolicy.FilterExpression = aero.ExpEq(aero.ExpStringBin("username"), aero.ExpStringVal(username))
	statement := aero.NewStatement(storage.AeroNsCache, AeroSetUserCache)
	recordSet, err := s.aero.Instance().Query(queryPolicy, statement)
	if err != nil && !err.Matches(types.KEY_NOT_FOUND_ERROR) {
		return nil, errors.ErrUserStorageGetCacheByUsername(err, ctx)
	}
	if recordSet == nil {
		return nil, err
	}
	// take the first result
	r := <-recordSet.Results()
	if r == nil {
		return nil, nil
	}
	if r.Err != nil {
		return nil, errors.ErrUserStorageGetCacheByUsername(r.Err, ctx)
	} else {
		return s.toUserCacheDomain(r.Record), nil
	}
}

func (s *UserStorageImpl) setCache(ctx context.Context, user *auth.User) error {
	s.l().Mth("set-cache").C(ctx).F(log.FF{"userId": user.Id}).Trc()
	key, err := aero.NewKey(storage.AeroNsCache, AeroSetUserCache, user.Id)
	if err != nil {
		return errors.ErrUserStorageAeroKey(err, ctx)
	}
	writePolicy := aero.NewWritePolicy(0, 3600)
	writePolicy.SendKey = true
	err = s.aero.Instance().Put(writePolicy, key, s.toUserCache(user))
	if err != nil {
		return errors.ErrUserStoragePutCache(err, ctx)
	}
	return nil
}

func (s *UserStorageImpl) CreateUser(ctx context.Context, user *auth.User) error {
	s.l().Mth("create").C(ctx).F(log.FF{"userId": user.Id}).Trc()
	dto := s.toUserDto(user)
	result := s.pg.Instance.Create(dto)
	if result.Error != nil {
		return errors.ErrUserStorageCreate(result.Error, ctx)
	}
	return nil
}

func (s *UserStorageImpl) UpdateUser(ctx context.Context, user *auth.User) error {
	l := s.l().Mth("update").C(ctx).F(log.FF{"userId": user.Id}).Trc()
	eg := goroutine.NewGroup(ctx).WithLogger(l)
	// save to DB
	eg.Go(func() error {
		dto := s.toUserDto(user)
		result := s.pg.Instance.Omit("created_at").Save(dto)
		if result.Error != nil {
			return errors.ErrUserStorageUpdate(result.Error, ctx)
		}
		return nil
	})
	// clear cache
	eg.Go(func() error {
		return s.clearCache(ctx, user.Id)
	})
	return eg.Wait()
}

func (s *UserStorageImpl) GetByUsername(ctx context.Context, username string) (*auth.User, error) {
	l := s.l().Mth("get").C(ctx).F(log.FF{"username": username}).Trc()
	if username == "" {
		return nil, nil
	}
	// check cache first
	usr, err := s.getFromCacheByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if usr != nil {
		l.Trc("found in cache")
		return usr, nil
	}
	// get from db
	dto := &user{}
	res := s.pg.Instance.Limit(1).Where("username = ?", username).Find(&dto)
	if res.Error != nil {
		return nil, errors.ErrUserStorageGetDb(res.Error, ctx)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	usr = s.toUserDomain(dto)
	// set cache
	err = s.setCache(ctx, usr)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (s *UserStorageImpl) GetUser(ctx context.Context, userId string) (*auth.User, error) {
	l := s.l().Mth("get").C(ctx).F(log.FF{"userId": userId}).Trc()
	if userId == "" {
		return nil, nil
	}
	// check cache first
	usr, err := s.getFromCacheById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if usr != nil {
		l.Trc("found in cache")
		return usr, nil
	}
	// get from db
	dto := &user{Id: userId}
	res := s.pg.Instance.Limit(1).Find(&dto)
	if res.Error != nil {
		return nil, errors.ErrUserStorageGetDb(res.Error, ctx)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	usr = s.toUserDomain(dto)
	// set cache
	err = s.setCache(ctx, usr)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (s *UserStorageImpl) GetUserByIds(ctx context.Context, userIds []string) ([]*auth.User, error) {
	s.l().Mth("get-ids").C(ctx).Trc()
	if len(userIds) == 0 {
		return []*auth.User{}, nil
	}
	var users []*user
	if err := s.pg.Instance.Find(&users, userIds).Error; err != nil {
		return nil, errors.ErrUserStorageGetByIds(err, ctx)
	}
	return s.toUsersDomain(users), nil
}

func (s *UserStorageImpl) DeleteUser(ctx context.Context, u *auth.User) error {
	l := s.l().C(ctx).Mth("delete").F(log.FF{"userId": u.Id}).Dbg()
	eg := goroutine.NewGroup(ctx).WithLogger(l)
	eg.Go(func() error {
		// save to DB
		if err := s.pg.Instance.Delete(&user{Id: u.Id}).Error; err != nil {
			return errors.ErrUserStorageDelete(err, ctx)
		}
		return nil
	})
	// clear cache
	eg.Go(func() error {
		return s.clearCache(ctx, u.Id)
	})
	return eg.Wait()
}
