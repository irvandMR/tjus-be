package setting

import (
	"time"

	"github.com/google/uuid"
)

type Setting struct {
	Id                uuid.UUID
	SiteName          string
	WaNumber          string
	WaButtonText      string
	WaMessageTemplate string
	UpdatedAt         time.Time
}

type SettingsRequest struct {
	SiteName          string `json:"site_name"`
	WaNumber          string `json:"wa_number"`
	WaButtonText      string `json:"wa_button_text"`
	WaMessageTemplate string `json:"wa_message_template"`
}

type UpdateSettingsRequest struct {
	Id uuid.UUID `json:"id"`
	SettingsRequest
}

type SettingsResponse struct {
	Id                uuid.UUID  `json:"id"`
	SiteName          string     `json:"site_name"`
	WaNumber          string     `json:"wa_number"`
	WaButtonText      string     `json:"wa_button_text"`
	WaMessageTemplate string     `json:"wa_message_template"`
	UpdatedAt         *time.Time `json:"updated_at"`
}
