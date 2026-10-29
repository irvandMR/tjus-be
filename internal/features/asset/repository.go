package asset

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssetRepository interface {
	// Transaction
	BeginTx(ctx context.Context) (pgx.Tx, error)

	// Asset
	InsertAsset(ctx context.Context, asset *Asset) (*Asset, error)
	InsertAssetTx(ctx context.Context, tx pgx.Tx, asset *Asset) (*Asset, error)
	UpdateAssetTx(ctx context.Context, tx pgx.Tx, asset *Asset) error
	UpdateAssetStatusTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, status string) error
	FindById(ctx context.Context, id uuid.UUID) (*Asset, error)
	FindBySlug(ctx context.Context, slug string) (*Asset, error)
	FindAll(ctx context.Context, filter AssetFilter) ([]Asset, int64, error)
	DeleteAsset(ctx context.Context, id uuid.UUID) error
	IncrementViewCount(ctx context.Context, id uuid.UUID) error
	// Asset Attribute
	BulkInsertAssetAttributes(ctx context.Context, attrs []AssetAttributes) error
	BulkInsertAssetAttributesTx(ctx context.Context, tx pgx.Tx, attrs []AssetAttributes) error
	DeleteAttributesByAssetIdTx(ctx context.Context, tx pgx.Tx, assetId uuid.UUID) error
	FindAttributesByAssetId(ctx context.Context, assetId uuid.UUID) ([]AssetAttributes, error)

	// Asset Media
	InsertAssetMedia(ctx context.Context, media *AssetMedia) (*AssetMedia, error)
	FindMediaByAssetId(ctx context.Context, assetId uuid.UUID) ([]AssetMedia, error)
	FindMediaById(ctx context.Context, id uuid.UUID) (*AssetMedia, error)
	DeleteMediaById(ctx context.Context, id uuid.UUID) error

	// Asset Owner
	InsertAssetOwner(ctx context.Context, owner *AssetOwner) error
	InsertAssetOwnerTx(ctx context.Context, tx pgx.Tx, owner *AssetOwner) error
	UpdateAssetOwnerTx(ctx context.Context, tx pgx.Tx, owner *AssetOwner) error
	FindOwnerByAssetId(ctx context.Context, assetId uuid.UUID) (*AssetOwner, error)
}

type assetRepository struct {
	db *pgxpool.Pool
}

func NewAssetRepository(db *pgxpool.Pool) AssetRepository {
	return &assetRepository{db: db}
}

// BeginTx membuka transaksi database secara langsung
func (ar *assetRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return ar.db.Begin(ctx)
}
