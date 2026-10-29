package setting

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type SettingsService struct {
	settingRepo SettingsRepository
}

func NewSettingsService(settingRepository SettingsRepository) *SettingsService {
	return &SettingsService{settingRepo: settingRepository}
}

func (ss *SettingsService) GetSetting(ctx context.Context) (*SettingsResponse, error) {
	res, err := ss.settingRepo.Find(ctx)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, errors.New("settings not found")
	}

	return &SettingsResponse{
		Id:                res.Id,
		SiteName:          res.SiteName,
		WaNumber:          res.WaNumber,
		WaButtonText:      res.WaButtonText,
		WaMessageTemplate: res.WaMessageTemplate,
	}, nil
}

func (ss *SettingsService) CreateSettings(ctx context.Context, request *SettingsRequest) error {

	err := ss.settingRepo.InsertSettings(ctx, Setting{
		SiteName:          request.SiteName,
		WaNumber:          request.WaNumber,
		WaButtonText:      request.WaButtonText,
		WaMessageTemplate: request.WaMessageTemplate,
	})

	return err
}

func (ss *SettingsService) EditSettings(ctx context.Context, request *UpdateSettingsRequest) (*SettingsResponse, error) {

	res, err := ss.validateSetting(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	set, err := ss.settingRepo.UpdateSettings(ctx, Setting{
		Id:                res.Id,
		SiteName:          request.SiteName,
		WaNumber:          request.WaNumber,
		WaButtonText:      request.WaButtonText,
		WaMessageTemplate: request.WaMessageTemplate,
	})

	if err != nil {
		return nil, err
	}

	return &SettingsResponse{
		Id:                res.Id,
		SiteName:          set.SiteName,
		WaNumber:          set.WaNumber,
		WaButtonText:      set.WaButtonText,
		WaMessageTemplate: set.WaMessageTemplate,
		UpdatedAt:         &set.UpdatedAt,
	}, nil

}

func (ss *SettingsService) validateSetting(ctx context.Context, id uuid.UUID) (*Setting, error) {
	settings, err := ss.settingRepo.FindById(ctx, id)

	if err != nil {
		return nil, err
	}
	if settings == nil {
		return nil, errors.New("settings not found")
	}
	return settings, nil
}
