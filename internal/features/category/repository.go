package category

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository interface {
	InsertCategory(ctx context.Context, request Category) (*Category, error)
	FindBySlug(ctx context.Context, slug string) (*Category, error)
	FindById(ctx context.Context, id uuid.UUID) (*Category, error)
	UpdateCategory(ctx context.Context, request Category) (*Category, error)
	FindAll(ctx context.Context) ([]Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error

	// Category Attributes
	InsertAttribute(ctx context.Context, attr CategoryAttributes) (*CategoryAttributes, error)
	BulkInsertAttributes(ctx context.Context, attrs []CategoryAttributes) ([]CategoryAttributes, error)
	FindAttributesByCategoryId(ctx context.Context, categoryId uuid.UUID) ([]CategoryAttributes, error)
	FindAttributeById(ctx context.Context, id uuid.UUID) (*CategoryAttributes, error)
	UpdateAttribute(ctx context.Context, attr CategoryAttributes) (*CategoryAttributes, error)
	DeleteAttribute(ctx context.Context, id uuid.UUID) error
}

type categoryRepository struct {
	db *pgxpool.Pool
}

// FindAll implements [CategoryRepository].
func (c *categoryRepository) FindAll(ctx context.Context) ([]Category, error) {
	query := `
		SELECT id, name, slug, sort_order, is_active, created_at, updated_at 
		FROM tbl_categories 
		ORDER BY sort_order ASC
	`
	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]Category, 0)
	for rows.Next() {
		var category Category
		err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.Slug,
			&category.SortOrder,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// UpdateCategory implements [CategoryRepository].
func (c *categoryRepository) UpdateCategory(ctx context.Context, request Category) (*Category, error) {
	dateNow := time.Now()
	query := `
		UPDATE tbl_categories
		SET 
			name = $1,
			slug = $2,
			sort_order = $3,
			is_active = $4,
			updated_at = $5
		WHERE id = $6
		RETURNING id, name, slug, sort_order, is_active, created_at, updated_at
	`
	var category Category
	err := c.db.QueryRow(ctx, query, request.Name, request.Slug, request.SortOrder, request.IsActive, dateNow, request.Id).Scan(&category.Id, &category.Name, &category.Slug, &category.SortOrder, &category.IsActive, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil

}

// FindById implements [CategoryRepository].
func (c *categoryRepository) FindById(ctx context.Context, id uuid.UUID) (*Category, error) {
	var category Category
	query := `
		SELECT id, name, slug, sort_order, is_active, created_at, updated_at 
		FROM tbl_categories 
		WHERE id = $1
	`
	err := c.db.QueryRow(ctx, query, id).Scan(&category.Id, &category.Name, &category.Slug, &category.SortOrder, &category.IsActive, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

// FindBySlug implements [CategoryRepository].
func (c *categoryRepository) FindBySlug(ctx context.Context, slug string) (*Category, error) {
	var category Category
	query := `
		SELECT id, name, slug, sort_order, is_active, created_at 
		FROM tbl_categories 
		WHERE slug = $1
	`
	err := c.db.QueryRow(ctx, query, slug).Scan(&category.Id, &category.Name, &category.Slug, &category.SortOrder, &category.IsActive, &category.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

// InsertCategory implements [CategoryRepository].
func (c *categoryRepository) InsertCategory(ctx context.Context, request Category) (*Category, error) {

	dateNow := time.Now()

	query := `
	INSERT INTO tbl_categories (name, slug, sort_order, is_active, created_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, name, slug, sort_order, is_active, created_at
	`
	var category Category
	err := c.db.QueryRow(ctx, query, request.Name, request.Slug, request.SortOrder, request.IsActive, dateNow).Scan(&category.Id, &category.Name, &category.Slug, &category.SortOrder, &category.IsActive, &category.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// DeleteCategory implements [CategoryRepository].
func (c *categoryRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tbl_categories WHERE id = $1`
	_, err := c.db.Exec(ctx, query, id)
	return err
}

// InsertAttribute implements [CategoryRepository].
func (c *categoryRepository) InsertAttribute(ctx context.Context, attr CategoryAttributes) (*CategoryAttributes, error) {
	dateNow := time.Now()
	query := `
		INSERT INTO tbl_catageries_attributes (
			category_id, name, data_type, unit, options, is_required, show_on_card, is_active, sort_order, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, category_id, name, data_type, unit, options, is_required, show_on_card, is_active, sort_order, created_at, updated_at
	`
	var result CategoryAttributes
	err := c.db.QueryRow(
		ctx,
		query,
		attr.CategoryId,
		attr.Name,
		attr.DataType,
		attr.Unit,
		attr.Options,
		attr.IsRequired,
		attr.ShowOnCard,
		attr.IsActive,
		attr.SortOrder,
		dateNow,
		dateNow,
	).Scan(
		&result.Id,
		&result.CategoryId,
		&result.Name,
		&result.DataType,
		&result.Unit,
		&result.Options,
		&result.IsRequired,
		&result.ShowOnCard,
		&result.IsActive,
		&result.SortOrder,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// BulkInsertAttributes menyimpan banyak atribut kategori sekaligus menggunakan pgx.Batch
func (c *categoryRepository) BulkInsertAttributes(ctx context.Context, attrs []CategoryAttributes) ([]CategoryAttributes, error) {
	if len(attrs) == 0 {
		return []CategoryAttributes{}, nil
	}

	dateNow := time.Now()
	query := `
		INSERT INTO tbl_catageries_attributes (
			category_id, name, data_type, unit, options, is_required, show_on_card, is_active, sort_order, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, category_id, name, data_type, unit, options, is_required, show_on_card, is_active, sort_order, created_at, updated_at
	`

	batch := &pgx.Batch{}
	for _, attr := range attrs {
		options := attr.Options
		if options == nil {
			options = []string{}
		}
		batch.Queue(
			query,
			attr.CategoryId,
			attr.Name,
			attr.DataType,
			attr.Unit,
			options,
			attr.IsRequired,
			attr.ShowOnCard,
			attr.IsActive,
			attr.SortOrder,
			dateNow,
			dateNow,
		)
	}

	br := c.db.SendBatch(ctx, batch)
	defer br.Close()

	results := make([]CategoryAttributes, 0, len(attrs))
	for range attrs {
		var result CategoryAttributes
		err := br.QueryRow().Scan(
			&result.Id,
			&result.CategoryId,
			&result.Name,
			&result.DataType,
			&result.Unit,
			&result.Options,
			&result.IsRequired,
			&result.ShowOnCard,
			&result.IsActive,
			&result.SortOrder,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

// FindAttributesByCategoryId implements [CategoryRepository].
func (c *categoryRepository) FindAttributesByCategoryId(ctx context.Context, categoryId uuid.UUID) ([]CategoryAttributes, error) {
	query := `
		SELECT id, category_id, name, data_type, unit, options, is_required, show_on_card, is_active, sort_order, created_at, updated_at
		FROM tbl_catageries_attributes
		WHERE category_id = $1
		ORDER BY sort_order ASC
	`
	rows, err := c.db.Query(ctx, query, categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attributes := make([]CategoryAttributes, 0)
	for rows.Next() {
		var attr CategoryAttributes
		err := rows.Scan(
			&attr.Id,
			&attr.CategoryId,
			&attr.Name,
			&attr.DataType,
			&attr.Unit,
			&attr.Options,
			&attr.IsRequired,
			&attr.ShowOnCard,
			&attr.IsActive,
			&attr.SortOrder,
			&attr.CreatedAt,
			&attr.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		attributes = append(attributes, attr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return attributes, nil
}

// FindAttributeById implements [CategoryRepository].
func (c *categoryRepository) FindAttributeById(ctx context.Context, id uuid.UUID) (*CategoryAttributes, error) {
	query := `
		SELECT id, category_id, name, data_type, unit, options, is_required, show_on_card, is_active, sort_order, created_at, updated_at
		FROM tbl_catageries_attributes
		WHERE id = $1
	`
	var attr CategoryAttributes
	err := c.db.QueryRow(ctx, query, id).Scan(
		&attr.Id,
		&attr.CategoryId,
		&attr.Name,
		&attr.DataType,
		&attr.Unit,
		&attr.Options,
		&attr.IsRequired,
		&attr.ShowOnCard,
		&attr.IsActive,
		&attr.SortOrder,
		&attr.CreatedAt,
		&attr.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &attr, nil
}

// UpdateAttribute implements [CategoryRepository].
func (c *categoryRepository) UpdateAttribute(ctx context.Context, attr CategoryAttributes) (*CategoryAttributes, error) {
	dateNow := time.Now()
	query := `
		UPDATE tbl_catageries_attributes
		SET 
			category_id = $1,
			name = $2,
			data_type = $3,
			unit = $4,
			options = $5,
			is_required = $6,
			show_on_card = $7,
			is_active = $8,
			sort_order = $9,
			updated_at = $10
		WHERE id = $11
		RETURNING id, category_id, name, data_type, unit, options, is_required, show_on_card, is_active, sort_order, created_at, updated_at
	`
	var result CategoryAttributes
	err := c.db.QueryRow(
		ctx,
		query,
		attr.CategoryId,
		attr.Name,
		attr.DataType,
		attr.Unit,
		attr.Options,
		attr.IsRequired,
		attr.ShowOnCard,
		attr.IsActive,
		attr.SortOrder,
		dateNow,
		attr.Id,
	).Scan(
		&result.Id,
		&result.CategoryId,
		&result.Name,
		&result.DataType,
		&result.Unit,
		&result.Options,
		&result.IsRequired,
		&result.ShowOnCard,
		&result.IsActive,
		&result.SortOrder,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteAttribute implements [CategoryRepository].
func (c *categoryRepository) DeleteAttribute(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tbl_catageries_attributes WHERE id = $1`
	_, err := c.db.Exec(ctx, query, id)
	return err
}

func NewCategoryRepository(db *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{db: db}
}
