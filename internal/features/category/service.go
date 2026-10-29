package category

import (
	"context"
	"errors"
	"strings"
	"tjus-be/internal/pkg/entity"

	"github.com/google/uuid"
)

type CategoryService struct {
	catRepo CategoryRepository
}

func NewCategoryService(catRepo CategoryRepository) *CategoryService {
	return &CategoryService{catRepo: catRepo}
}

func (c *CategoryService) CreateCategory(ctx context.Context, req *CreatedCatRequest) (*CategoryResponse, error) {
	converSlug := generateSlug(req.Name)
	existingCat, err := c.catRepo.FindBySlug(ctx, converSlug)
	if err != nil {
		return nil, err
	}

	if existingCat != nil {
		return nil, errors.New("category already exists")
	}

	cat, err := c.catRepo.InsertCategory(ctx, Category{
		Name:      req.Name,
		Slug:      converSlug,
		SortOrder: req.SortOrder,
		IsActive:  req.IsActive,
	})

	if err != nil {
		return nil, err
	}

	return &CategoryResponse{
		Id:        cat.Id,
		Name:      cat.Name,
		Slug:      cat.Slug,
		SortOrder: cat.SortOrder,
		IsActive:  cat.IsActive,
	}, nil
}

func (c *CategoryService) EditCategory(ctx context.Context, req *UpdateCatRequest) (*CategoryResponse, error) {
	// 1. Validasi apakah kategori yang mau diedit ada di database
	if _, err := c.validateCategory(ctx, req.Id); err != nil {
		return nil, err
	}

	// 2. Validasi apakah slug baru bentrok dengan data milik kategori lain
	converSlug := generateSlug(req.Name)
	if err := c.validateSlug(ctx, converSlug, req.Id); err != nil {
		return nil, err
	}

	// 3. Eksekusi update ke database
	cat, err := c.catRepo.UpdateCategory(ctx, Category{
		BaseEntity: entity.BaseEntity{Id: req.Id},
		Name:       req.Name,
		Slug:       converSlug,
		SortOrder:  req.SortOrder,
		IsActive:   req.IsActive,
	})
	if err != nil {
		return nil, err
	}

	return &CategoryResponse{
		Id:        cat.Id,
		Name:      cat.Name,
		Slug:      cat.Slug,
		SortOrder: cat.SortOrder,
		IsActive:  cat.IsActive,
	}, nil
}

func (c *CategoryService) GetAllCategories(ctx context.Context) ([]CategoryResponse, error) {
	categories, err := c.catRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]CategoryResponse, 0, len(categories))
	for _, cat := range categories {
		responses = append(responses, CategoryResponse{
			Id:        cat.Id,
			Name:      cat.Name,
			Slug:      cat.Slug,
			SortOrder: cat.SortOrder,
			IsActive:  cat.IsActive,
		})
	}

	return responses, nil
}

func (c *CategoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	// 1. Validasi apakah kategori yang mau dihapus ada
	if _, err := c.validateCategory(ctx, id); err != nil {
		return err
	}

	// 2. Eksekusi hapus di repository
	return c.catRepo.DeleteCategory(ctx, id)
}

func (c *CategoryService) CreateAttribute(ctx context.Context, req *CreatedCatAttrRequest) ([]CategoryAttributeResponse, error) {
	// 1. Validasi apakah category_id valid dan ada di database
	if _, err := c.validateCategory(ctx, req.CategoryId); err != nil {
		return nil, err
	}

	// 2. Siapkan slice entity untuk batch insert
	attrEntities := make([]CategoryAttributes, 0, len(req.Attributes))
	for _, item := range req.Attributes {
		options := item.Options
		if options == nil {
			options = []string{}
		}
		attrEntities = append(attrEntities, CategoryAttributes{
			CategoryId: req.CategoryId,
			Name:       item.Name,
			DataType:   item.DataType,
			Unit:       item.Unit,
			Options:    options,
			IsRequired: item.IsRequired,
			ShowOnCard: item.ShowOnCard,
			IsActive:   item.IsActive,
			SortOrder:  item.SortOrder,
		})
	}

	// 3. Simpan data sekaligus ke database menggunakan batch insert
	savedAttrs, err := c.catRepo.BulkInsertAttributes(ctx, attrEntities)
	if err != nil {
		return nil, err
	}

	// 4. Susun response
	responses := make([]CategoryAttributeResponse, 0, len(savedAttrs))
	for _, a := range savedAttrs {
		responses = append(responses, CategoryAttributeResponse{
			Id:         a.Id,
			CategoryId: a.CategoryId,
			Name:       a.Name,
			DataType:   a.DataType,
			Unit:       a.Unit,
			Options:    a.Options,
			IsRequired: a.IsRequired,
			ShowOnCard: a.ShowOnCard,
			IsActive:   a.IsActive,
			SortOrder:  a.SortOrder,
		})
	}

	return responses, nil
}

func (c *CategoryService) GetAttributesByCategoryId(ctx context.Context, categoryId uuid.UUID) ([]CategoryAttributeResponse, error) {
	// Validasi apakah kategori ada
	if _, err := c.validateCategory(ctx, categoryId); err != nil {
		return nil, err
	}

	attrs, err := c.catRepo.FindAttributesByCategoryId(ctx, categoryId)
	if err != nil {
		return nil, err
	}

	responses := make([]CategoryAttributeResponse, 0, len(attrs))
	for _, a := range attrs {
		responses = append(responses, CategoryAttributeResponse{
			Id:         a.Id,
			CategoryId: a.CategoryId,
			Name:       a.Name,
			DataType:   a.DataType,
			Unit:       a.Unit,
			Options:    a.Options,
			IsRequired: a.IsRequired,
			ShowOnCard: a.ShowOnCard,
			IsActive:   a.IsActive,
			SortOrder:  a.SortOrder,
		})
	}

	return responses, nil
}

func (c *CategoryService) EditAttribute(ctx context.Context, req *UpdateCatAttrRequest) (*CategoryAttributeResponse, error) {
	// 1. Validasi apakah attribute ada
	existingAttr, err := c.catRepo.FindAttributeById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if existingAttr == nil {
		return nil, errors.New("attribute not found")
	}

	// 2. Validasi apakah category_id ada
	if _, err := c.validateCategory(ctx, req.CategoryId); err != nil {
		return nil, err
	}

	// 3. Update attribute
	attr, err := c.catRepo.UpdateAttribute(ctx, CategoryAttributes{
		BaseEntity: entity.BaseEntity{Id: req.Id},
		CategoryId: req.CategoryId,
		Name:       req.Name,
		DataType:   req.DataType,
		Unit:       req.Unit,
		Options:    req.Options,
		IsRequired: req.IsRequired,
		ShowOnCard: req.ShowOnCard,
		IsActive:   req.IsActive,
		SortOrder:  req.SortOrder,
	})
	if err != nil {
		return nil, err
	}

	return &CategoryAttributeResponse{
		Id:         attr.Id,
		CategoryId: attr.CategoryId,
		Name:       attr.Name,
		DataType:   attr.DataType,
		Unit:       attr.Unit,
		Options:    attr.Options,
		IsRequired: attr.IsRequired,
		ShowOnCard: attr.ShowOnCard,
		IsActive:   attr.IsActive,
		SortOrder:  attr.SortOrder,
	}, nil
}

func (c *CategoryService) DeleteAttribute(ctx context.Context, id uuid.UUID) error {
	// 1. Validasi apakah attribute ada
	existingAttr, err := c.catRepo.FindAttributeById(ctx, id)
	if err != nil {
		return err
	}
	if existingAttr == nil {
		return errors.New("attribute not found")
	}

	// 2. Hapus
	return c.catRepo.DeleteAttribute(ctx, id)
}

func generateSlug(text string) string {
	lower := strings.ToLower(text)
	return strings.ReplaceAll(lower, " ", "-")
}

func (c *CategoryService) validateCategory(ctx context.Context, id uuid.UUID) (*Category, error) {
	category, err := c.catRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}
	return category, nil
}

func (c *CategoryService) validateSlug(ctx context.Context, slug string, currentId uuid.UUID) error {
	existingCat, err := c.catRepo.FindBySlug(ctx, slug)
	if err != nil {
		return err
	}
	if existingCat != nil && existingCat.Id != currentId {
		return errors.New("category with this slug already exists")
	}
	return nil
}
