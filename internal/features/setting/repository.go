package setting

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SettingsRepository interface {
	InsertSettings(ctx context.Context, request Setting) error
	UpdateSettings(ctx context.Context, request Setting) (*Setting, error)
	FindById(ctx context.Context, id uuid.UUID) (*Setting, error)
	Find(ctx context.Context) (*Setting, error)
}

type settingsRepository struct {
	db *pgxpool.Pool
}

func (sr *settingsRepository) InsertSettings(ctx context.Context, request Setting) error {

	query := `
		INSERT INTO tbl_settings (site_name, wa_number, wa_button_text, wa_message_template)
		VALUES ($1, $2, $3, $4)
	`

	_, err := sr.db.Exec(ctx, query, request.SiteName, request.WaNumber, request.WaButtonText, request.WaMessageTemplate)

	return err
}

func (sr *settingsRepository) UpdateSettings(ctx context.Context, request Setting) (*Setting, error) {
	dateTime := time.Now()

	query := `
		UPDATE tbl_settings
		SET 
			site_name = $1,
			wa_number = $2,
			wa_button_text = $3,
			wa_message_template = $4,
			updated_at = $5
		WHERE id = $6
		RETURNING id, site_name, wa_number, wa_button_text, wa_message_template, updated_at
	`

	var setting Setting
	err := sr.db.QueryRow(
		ctx,
		query,
		request.SiteName,
		request.WaNumber,
		request.WaButtonText,
		request.WaMessageTemplate,
		dateTime,
		request.Id,
	).Scan(
		&setting.Id,
		&setting.SiteName,
		&setting.WaNumber,
		&setting.WaButtonText,
		&setting.WaMessageTemplate,
		&setting.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &setting, nil
}

func (sr *settingsRepository) FindById(ctx context.Context, id uuid.UUID) (*Setting, error) {
	var setting Setting
	query := `
		SELECT id, site_name, wa_number, wa_button_text, wa_message_template, updated_at
		FROM tbl_settings
		WHERE id = $1
	`

	err := sr.db.QueryRow(ctx, query, id).Scan(&setting.Id, &setting.SiteName, &setting.WaNumber, &setting.WaButtonText, &setting.WaMessageTemplate, &setting.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	return &setting, nil
}

func (sr *settingsRepository) Find(ctx context.Context) (*Setting, error) {
	var setting Setting
	query := `
		SELECT id, site_name, wa_number, wa_button_text, wa_message_template, updated_at
		FROM tbl_settings LIMIT 1
	`

	err := sr.db.QueryRow(ctx, query).Scan(&setting.Id, &setting.SiteName, &setting.WaNumber, &setting.WaButtonText, &setting.WaMessageTemplate, &setting.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	return &setting, nil
}

func NewSettingsRepository(db *pgxpool.Pool) SettingsRepository {
	return &settingsRepository{db: db}
}
