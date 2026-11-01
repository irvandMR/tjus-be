package category

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// mockCategoryRepository memisahkan seluruh kode pemeran pengganti (Mock)
// agar file service_test.go tetap ringkas dan fokus pada skenario pengujian.
type mockCategoryRepository struct {
	mock.Mock
}

func (m *mockCategoryRepository) InsertCategory(ctx context.Context, request Category) (*Category, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Category), args.Error(1)
}

func (m *mockCategoryRepository) FindBySlug(ctx context.Context, slug string) (*Category, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Category), args.Error(1)
}

func (m *mockCategoryRepository) FindById(ctx context.Context, id uuid.UUID) (*Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Category), args.Error(1)
}

func (m *mockCategoryRepository) UpdateCategory(ctx context.Context, request Category) (*Category, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Category), args.Error(1)
}

func (m *mockCategoryRepository) FindAll(ctx context.Context) ([]Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Category), args.Error(1)
}

func (m *mockCategoryRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockCategoryRepository) InsertAttribute(ctx context.Context, attr CategoryAttributes) (*CategoryAttributes, error) {
	args := m.Called(ctx, attr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CategoryAttributes), args.Error(1)
}

func (m *mockCategoryRepository) BulkInsertAttributes(ctx context.Context, attrs []CategoryAttributes) ([]CategoryAttributes, error) {
	args := m.Called(ctx, attrs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]CategoryAttributes), args.Error(1)
}

func (m *mockCategoryRepository) FindAttributesByCategoryId(ctx context.Context, categoryId uuid.UUID) ([]CategoryAttributes, error) {
	args := m.Called(ctx, categoryId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]CategoryAttributes), args.Error(1)
}

func (m *mockCategoryRepository) FindAttributeById(ctx context.Context, id uuid.UUID) (*CategoryAttributes, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CategoryAttributes), args.Error(1)
}

func (m *mockCategoryRepository) UpdateAttribute(ctx context.Context, attr CategoryAttributes) (*CategoryAttributes, error) {
	args := m.Called(ctx, attr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CategoryAttributes), args.Error(1)
}

func (m *mockCategoryRepository) DeleteAttribute(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

