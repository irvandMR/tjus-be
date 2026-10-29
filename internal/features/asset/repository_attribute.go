package asset

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const insertAttributeQuery = `
	INSERT INTO tbl_asset_attribute_values (
		asset_id, attribute_id, custome_label, value_text, value_number, sort_order
	) VALUES ($1, $2, $3, $4, $5, $6)
`

// BulkInsertAssetAttributes untuk bulk query normal tanpa transaksi
func (ar *assetRepository) BulkInsertAssetAttributes(ctx context.Context, attrs []AssetAttributes) error {
	if len(attrs) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, attr := range attrs {
		batch.Queue(insertAttributeQuery, attr.AssetId, attr.CategoryAttributeId, attr.CustomeLabel, attr.ValueText, attr.ValueNumber, attr.SortOrder)
	}

	br := ar.db.SendBatch(ctx, batch)
	defer br.Close()

	for range attrs {
		_, err := br.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

// BulkInsertAssetAttributesTx untuk bulk query di dalam transaksi database yang dioper eksplisit
func (ar *assetRepository) BulkInsertAssetAttributesTx(ctx context.Context, tx pgx.Tx, attrs []AssetAttributes) error {
	if len(attrs) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, attr := range attrs {
		batch.Queue(insertAttributeQuery, attr.AssetId, attr.CategoryAttributeId, attr.CustomeLabel, attr.ValueText, attr.ValueNumber, attr.SortOrder)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	for range attrs {
		_, err := br.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

func (ar *assetRepository) DeleteAttributesByAssetIdTx(ctx context.Context, tx pgx.Tx, assetId uuid.UUID) error {
	query := `DELETE FROM tbl_asset_attribute_values WHERE asset_id = $1`
	_, err := tx.Exec(ctx, query, assetId)
	return err
}

func (ar *assetRepository) FindAttributesByAssetId(ctx context.Context, assetId uuid.UUID) ([]AssetAttributes, error) {
	query := `
		SELECT 
			aav.id,
			aav.asset_id,
			aav.attribute_id,
			ca.name,
			ca.data_type,
			ca.unit,
			aav.custome_label,
			aav.value_text,
			aav.value_number,
			aav.sort_order
		FROM tbl_asset_attribute_values aav
		LEFT JOIN tbl_catageries_attributes ca ON ca.id = aav.attribute_id
		WHERE aav.asset_id = $1
		ORDER BY aav.sort_order ASC
	`

	rows, err := ar.db.Query(ctx, query, assetId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]AssetAttributes, 0)
	for rows.Next() {
		var item AssetAttributes
		err := rows.Scan(
			&item.Id,
			&item.AssetId,
			&item.CategoryAttributeId,
			&item.Name,
			&item.DataType,
			&item.Unit,
			&item.CustomeLabel,
			&item.ValueText,
			&item.ValueNumber,
			&item.SortOrder,
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
