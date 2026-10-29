package asset

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const insertAssetQuery = `
	INSERT INTO tbl_assets (
		code,
		slug,
		category_id,
		title,
		location,
		description,
		badge,
		price_label,
		price,
		price_negotiable,
		view_count,
		status,
		created_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
	)
	RETURNING 
		id,
		code,
		slug,
		category_id,
		title,
		location,
		description,
		badge,
		price_label,
		price,
		price_negotiable,
		view_count,
		status,
		created_at,
		updated_at
`

// InsertAsset untuk query normal tanpa transaksi
func (ar *assetRepository) InsertAsset(ctx context.Context, asset *Asset) (*Asset, error) {
	dateNow := time.Now()
	var result Asset
	err := ar.db.QueryRow(
		ctx,
		insertAssetQuery,
		asset.Code,
		asset.Slug,
		asset.CategoryId,
		asset.Title,
		asset.Location,
		asset.Description,
		asset.Badge,
		asset.PriceLabel,
		asset.Price,
		asset.PriceNegotiable,
		asset.ViewCount,
		asset.Status,
		dateNow,
	).Scan(
		&result.Id,
		&result.Code,
		&result.Slug,
		&result.CategoryId,
		&result.Title,
		&result.Location,
		&result.Description,
		&result.Badge,
		&result.PriceLabel,
		&result.Price,
		&result.PriceNegotiable,
		&result.ViewCount,
		&result.Status,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// InsertAssetTx untuk query di dalam transaksi database yang dioper eksplisit
func (ar *assetRepository) InsertAssetTx(ctx context.Context, tx pgx.Tx, asset *Asset) (*Asset, error) {
	dateNow := time.Now()
	var result Asset
	err := tx.QueryRow(
		ctx,
		insertAssetQuery,
		asset.Code,
		asset.Slug,
		asset.CategoryId,
		asset.Title,
		asset.Location,
		asset.Description,
		asset.Badge,
		asset.PriceLabel,
		asset.Price,
		asset.PriceNegotiable,
		asset.ViewCount,
		asset.Status,
		dateNow,
	).Scan(
		&result.Id,
		&result.Code,
		&result.Slug,
		&result.CategoryId,
		&result.Title,
		&result.Location,
		&result.Description,
		&result.Badge,
		&result.PriceLabel,
		&result.Price,
		&result.PriceNegotiable,
		&result.ViewCount,
		&result.Status,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (ar *assetRepository) FindById(ctx context.Context, id uuid.UUID) (*Asset, error) {
	var asset Asset
	query := `
		SELECT 
			id,
			code,
			slug,
			category_id,
			title,
			location,
			description,
			badge,
			price_label,
			price,
			price_negotiable,
			view_count,
			status,
			created_at,
			updated_at
		FROM tbl_assets
		WHERE id = $1
	`

	err := ar.db.QueryRow(ctx, query, id).Scan(
		&asset.Id,
		&asset.Code,
		&asset.Slug,
		&asset.CategoryId,
		&asset.Title,
		&asset.Location,
		&asset.Description,
		&asset.Badge,
		&asset.PriceLabel,
		&asset.Price,
		&asset.PriceNegotiable,
		&asset.ViewCount,
		&asset.Status,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &asset, nil
}

func (ar *assetRepository) FindBySlug(ctx context.Context, slug string) (*Asset, error) {
	var asset Asset
	query := `
		SELECT 
			id,
			code,
			slug,
			category_id,
			title,
			location,
			description,
			badge,
			price_label,
			price,
			price_negotiable,
			view_count,
			status,
			created_at,
			updated_at
		FROM tbl_assets
		WHERE slug = $1
	`

	err := ar.db.QueryRow(ctx, query, slug).Scan(
		&asset.Id,
		&asset.Code,
		&asset.Slug,
		&asset.CategoryId,
		&asset.Title,
		&asset.Location,
		&asset.Description,
		&asset.Badge,
		&asset.PriceLabel,
		&asset.Price,
		&asset.PriceNegotiable,
		&asset.ViewCount,
		&asset.Status,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &asset, nil
}

func (ar *assetRepository) FindAll(ctx context.Context, filter AssetFilter) ([]Asset, int64, error) {
	var whereClauses []string
	var args []any
	argIdx := 1

	// 1. Status filter (Default: ACTIVE jika nil/kosong)
	if filter.Status != nil && *filter.Status != "" {
		if strings.ToUpper(*filter.Status) != "ALL" {
			whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argIdx))
			args = append(args, *filter.Status)
			argIdx++
		}
	} else {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, "ACTIVE")
		argIdx++
	}

	// 2. Category filter
	if filter.CategoryId != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("category_id = $%d", argIdx))
		args = append(args, *filter.CategoryId)
		argIdx++
	}

	// 3. Search query (title, location, code)
	if strings.TrimSpace(filter.Search) != "" {
		keyword := "%" + strings.TrimSpace(filter.Search) + "%"
		whereClauses = append(whereClauses, fmt.Sprintf("(title ILIKE $%d OR location ILIKE $%d OR code ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, keyword)
		argIdx++
	}

	// 4. Min Price filter
	if filter.MinPrice != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("price >= $%d", argIdx))
		args = append(args, *filter.MinPrice)
		argIdx++
	}

	// 5. Max Price filter
	if filter.MaxPrice != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("price <= $%d", argIdx))
		args = append(args, *filter.MaxPrice)
		argIdx++
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Hitung total data untuk pagination meta
	var total int64
	countQuery := "SELECT COUNT(*) FROM tbl_assets" + whereSQL
	if err := ar.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Ambil data aset sesuai limit dan offset
	query := fmt.Sprintf(`
		SELECT 
			id,
			code,
			slug,
			category_id,
			title,
			location,
			description,
			badge,
			price_label,
			price,
			price_negotiable,
			view_count,
			status,
			created_at,
			updated_at
		FROM tbl_assets
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereSQL, argIdx, argIdx+1)

	queryArgs := append(args, filter.Limit, filter.Offset)
	rows, err := ar.db.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	assets := make([]Asset, 0)
	for rows.Next() {
		var a Asset
		err := rows.Scan(
			&a.Id,
			&a.Code,
			&a.Slug,
			&a.CategoryId,
			&a.Title,
			&a.Location,
			&a.Description,
			&a.Badge,
			&a.PriceLabel,
			&a.Price,
			&a.PriceNegotiable,
			&a.ViewCount,
			&a.Status,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		assets = append(assets, a)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return assets, total, nil
}

func (ar *assetRepository) UpdateAssetTx(ctx context.Context, tx pgx.Tx, asset *Asset) error {
	dateNow := time.Now()
	query := `
		UPDATE tbl_assets SET
			title = $1,
			slug = $2,
			location = $3,
			description = $4,
			badge = $5,
			price_label = $6,
			price = $7,
			price_negotiable = $8,
			updated_at = $9
		WHERE id = $10
	`
	_, err := tx.Exec(
		ctx,
		query,
		asset.Title,
		asset.Slug,
		asset.Location,
		asset.Description,
		asset.Badge,
		asset.PriceLabel,
		asset.Price,
		asset.PriceNegotiable,
		dateNow,
		asset.Id,
	)
	return err
}

func (ar *assetRepository) UpdateAssetStatusTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, status string) error {
	dateNow := time.Now()
	query := `
		UPDATE tbl_assets SET
			status = $1,
			updated_at = $2
		WHERE id = $3
	`
	_, err := tx.Exec(ctx, query, status, dateNow, id)
	return err
}

func (ar *assetRepository) DeleteAsset(ctx context.Context, id uuid.UUID) error {
	dateNow := time.Now()
	query := `
		UPDATE tbl_assets SET
			status = 'DELETED',
			updated_at = $1
		WHERE id = $2
	`
	_, err := ar.db.Exec(ctx, query, dateNow, id)
	return err
}

func (ar *assetRepository) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE tbl_assets SET
			view_count = view_count + 1
		WHERE id = $1
	`
	_, err := ar.db.Exec(ctx, query, id)
	return err
}
