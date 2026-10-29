package asset

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (ar *assetRepository) InsertAssetMedia(ctx context.Context, media *AssetMedia) (*AssetMedia, error) {
	dateNow := time.Now()

	query := `
		INSERT INTO tbl_asset_media (
			asset_id,
			type_media,
			file_path,
			sort_order,
			created_at
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, asset_id, type_media, file_path, sort_order, created_at
	`

	var result AssetMedia
	err := ar.db.QueryRow(
		ctx,
		query,
		media.AssetId,
		media.TypeMedia,
		media.FilePath,
		media.SortOrder,
		dateNow,
	).Scan(
		&result.Id,
		&result.AssetId,
		&result.TypeMedia,
		&result.FilePath,
		&result.SortOrder,
		&result.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ar *assetRepository) FindMediaByAssetId(ctx context.Context, assetId uuid.UUID) ([]AssetMedia, error) {
	query := `
		SELECT 
			id,
			asset_id,
			type_media,
			file_path,
			sort_order,
			created_at
		FROM tbl_asset_media
		WHERE asset_id = $1
		ORDER BY sort_order ASC
	`

	rows, err := ar.db.Query(ctx, query, assetId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]AssetMedia, 0)
	for rows.Next() {
		var item AssetMedia
		err := rows.Scan(
			&item.Id,
			&item.AssetId,
			&item.TypeMedia,
			&item.FilePath,
			&item.SortOrder,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (ar *assetRepository) FindMediaById(ctx context.Context, id uuid.UUID) (*AssetMedia, error) {
	query := `
		SELECT 
			id,
			asset_id,
			type_media,
			file_path,
			sort_order,
			created_at
		FROM tbl_asset_media
		WHERE id = $1
	`
	var item AssetMedia
	err := ar.db.QueryRow(ctx, query, id).Scan(
		&item.Id,
		&item.AssetId,
		&item.TypeMedia,
		&item.FilePath,
		&item.SortOrder,
		&item.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (ar *assetRepository) DeleteMediaById(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tbl_asset_media WHERE id = $1`
	_, err := ar.db.Exec(ctx, query, id)
	return err
}
