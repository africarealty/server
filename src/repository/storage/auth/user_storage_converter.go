package auth

import (
	"encoding/json"
	aero "github.com/aerospike/aerospike-client-go/v6"
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/kit/auth"
	"github.com/africarealty/server/src/kit/storages/pg"
)

func (s *UserStorageImpl) toUserDto(u *domain.User) *user {
	if u == nil {
		return nil
	}
	dto := &user{
		Id:          u.Id,
		Username:    u.Username,
		Password:    pg.StringToNull(u.Password),
		Type:        u.Type,
		ActivatedAt: u.ActivatedAt,
		LockedAt:    u.LockedAt,
	}
	det := &userDetails{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Groups:    u.Groups,
		Roles:     u.Roles,
		Owner:     u.Owner,
		Agent:     u.Agent,
	}
	var detailsBytes []byte
	detailsBytes, _ = json.Marshal(det)
	dto.Details = string(detailsBytes)
	return dto
}

func (s *UserStorageImpl) toUserCacheDomain(rec *aero.Record) *domain.User {
	if rec == nil {
		return nil
	}
	body := rec.Bins["user"].(string)
	user := &domain.User{}
	_ = json.Unmarshal([]byte(body), user)
	return user
}

func (s *UserStorageImpl) toUserCache(user *domain.User) aero.BinMap {
	usrBytes, _ := json.Marshal(user)
	return aero.BinMap{
		"username": user.Username,
		"user":     string(usrBytes),
	}
}

func (s *UserStorageImpl) toUserDomain(dto *user) *domain.User {
	if dto == nil {
		return nil
	}
	det := &userDetails{}
	_ = json.Unmarshal([]byte(dto.Details), det)
	return &domain.User{
		User: auth.User{
			Id:          dto.Id,
			Username:    dto.Username,
			Password:    pg.NullToString(dto.Password),
			Type:        dto.Type,
			FirstName:   det.FirstName,
			LastName:    det.LastName,
			ActivatedAt: dto.ActivatedAt,
			LockedAt:    dto.LockedAt,
			Groups:      det.Groups,
			Roles:       det.Roles,
		},
		Owner: det.Owner,
		Agent: det.Agent,
	}
}

func (s *UserStorageImpl) toUsersDomain(dtos []*user) []*domain.User {
	var res []*domain.User
	for _, d := range dtos {
		res = append(res, s.toUserDomain(d))
	}
	return res
}
