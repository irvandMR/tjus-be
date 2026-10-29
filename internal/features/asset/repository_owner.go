package asset

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const insertOwnerQuery = `
	INSERT INTO tbl_asset_owner_info (
		asset_id,
		owner_name,
		owner_contact,
		internal_notes,
		updated_at
	) VALUES ($1, $2, $3, $4, $5)
`

// InsertAssetOwner untuk query normal tanpa transaksi
func (ar *assetRepository) InsertAssetOwner(ctx context.Context, owner *AssetOwner) error {
	dateNow := time.Now()
	_, err := ar.db.Exec(
		ctx,
		insertOwnerQuery,
		owner.AssetId,
		owner.OwnerName,
		owner.OwnerContact,
		owner.InternalNotes,
		dateNow,
	)
	return err
}

// InsertAssetOwnerTx untuk query di dalam transaksi database yang dioper eksplisit
func (ar *assetRepository) InsertAssetOwnerTx(ctx context.Context, tx pgx.Tx, owner *AssetOwner) error {
	dateNow := time.Now()
	_, err := tx.Exec(
		ctx,
		insertOwnerQuery,
		owner.AssetId,
		owner.OwnerName,
		owner.OwnerContact,
		owner.InternalNotes,
		dateNow,
	)
	return err
}

func (ar *assetRepository) UpdateAssetOwnerTx(ctx context.Context, tx pgx.Tx, owner *AssetOwner) error {
	dateNow := time.Now()
	query := `
		UPDATE tbl_asset_owner_info SET
			owner_name = $1,
			owner_contact = $2,
			internal_notes = $3,
			updated_at = $4
		WHERE asset_id = $5
	`
	_, err := tx.Exec(
		ctx,
		query,
		owner.OwnerName,
		owner.OwnerContact,
		owner.InternalNotes,
		dateNow,
		owner.AssetId,
	)
	return err
}

func (ar *assetRepository) FindOwnerByAssetId(ctx context.Context, assetId uuid.UUID) (*AssetOwner, error) {
	var owner AssetOwner

	query := `
		SELECT 
			asset_id,
			owner_name,
			owner_contact,
			internal_notes,
			updated_at
		FROM tbl_asset_owner_info
		WHERE asset_id = $1
	`

	err := ar.db.QueryRow(ctx, query, assetId).Scan(
		&owner.AssetId,
		&owner.OwnerName,
		&owner.OwnerContact,
		&owner.InternalNotes,
		&owner.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &owner, nil
}
