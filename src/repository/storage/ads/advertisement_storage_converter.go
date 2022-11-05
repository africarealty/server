package ads

import (
	"encoding/json"
	aero "github.com/aerospike/aerospike-client-go/v6"
	"github.com/africarealty/server/src/domain"
)

func (s *AdvertisementStorageImpl) toAdsDto(a *domain.Advertisement) *advertisement {
	if a == nil {
		return nil
	}
	dto := &advertisement{
		Id:          a.Id,
		Code:        a.Code,
		UserId:      a.Details.Participants.Advertiser.UserId,
		Status:      a.Status,
		SubStatus:   a.SubStatus,
		Type:        a.Type,
		SubType:     a.SubType,
		ActivatedAt: a.ActivatedAt,
		ClosedAt:    a.ClosedAt,
	}
	var detailsBytes []byte
	detailsBytes, _ = json.Marshal(a.Details)
	dto.Details = string(detailsBytes)
	return dto
}

func (s *AdvertisementStorageImpl) toAdsCacheDomain(rec *aero.Record) *domain.Advertisement {
	if rec == nil {
		return nil
	}
	body := rec.Bins["ads"].(string)
	ads := &domain.Advertisement{}
	_ = json.Unmarshal([]byte(body), ads)
	return ads
}

func (s *AdvertisementStorageImpl) toAdsCache(ads *domain.Advertisement) aero.BinMap {
	adsBytes, _ := json.Marshal(ads)
	return aero.BinMap{
		"user_id": ads.Details.Participants.Advertiser.UserId,
		"ads":     string(adsBytes),
	}
}

func (s *AdvertisementStorageImpl) toAdsDomain(dto *advertisement) *domain.Advertisement {
	if dto == nil {
		return nil
	}
	det := &domain.AdDetails{}
	_ = json.Unmarshal([]byte(dto.Details), det)
	return &domain.Advertisement{
		Id:          dto.Id,
		Code:        dto.Code,
		Status:      dto.Status,
		SubStatus:   dto.SubStatus,
		Type:        dto.Type,
		SubType:     dto.SubType,
		Details:     det,
		ActivatedAt: dto.ActivatedAt,
		ClosedAt:    dto.ClosedAt,
	}
}

func (s *AdvertisementStorageImpl) toAdsIndex(ads *domain.Advertisement) *esIdxAds {
	if ads == nil {
		return nil
	}
	return &esIdxAds{
		Code:      ads.Code,
		UserId:    ads.Details.Participants.Advertiser.UserId,
		Status:    ads.Status,
		SubStatus: ads.SubStatus,
		Type:      ads.Type,
		SubType:   ads.SubType,
	}
}
