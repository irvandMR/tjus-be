package category

import (
	"context"
	"testing"
	"tjus-be/internal/pkg/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCategoryService_CreateCategory(t *testing.T) {
	catId := uuid.New()

	tests := []struct {
		name        string
		req         *CreatedCatRequest
		setupMock   func(m *mockCategoryRepository)
		wantErr     bool
		expectedErr string
	}{
		{
			name: "sukses membuat kategori baru dengan slug otomatis",
			req: &CreatedCatRequest{
				Name:      "Alat Berat",
				SortOrder: 1,
				IsActive:  true,
			},
			setupMock: func(m *mockCategoryRepository) {
				// Langkah 1: Cek slug belum ada di DB (nil, nil)
				m.On("FindBySlug", mock.Anything, "alat-berat").Return(nil, nil)
				// Langkah 2: Simpan ke DB berhasil
				m.On("InsertCategory", mock.Anything, Category{
					Name:      "Alat Berat",
					Slug:      "alat-berat",
					SortOrder: 1,
					IsActive:  true,
				}).Return(&Category{
					BaseEntity: entity.BaseEntity{Id: catId},
					Name:       "Alat Berat",
					Slug:       "alat-berat",
					SortOrder:  1,
					IsActive:   true,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "gagal karena slug kategori sudah terdaftar",
			req: &CreatedCatRequest{
				Name:      "Mobil Bekas",
				SortOrder: 2,
				IsActive:  true,
			},
			setupMock: func(m *mockCategoryRepository) {
				// Slug sudah ada -> InsertCategory TIDAK BOLEH dipanggil sama sekali!
				m.On("FindBySlug", mock.Anything, "mobil-bekas").Return(&Category{
					BaseEntity: entity.BaseEntity{Id: catId},
					Name:       "Mobil Bekas",
					Slug:       "mobil-bekas",
				}, nil)
			},
			wantErr:     true,
			expectedErr: "category already exists",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mockCategoryRepository)
			tc.setupMock(mockRepo)

			svc := NewCategoryService(mockRepo)
			resp, err := svc.CreateCategory(context.Background(), tc.req)

			if tc.wantErr {
				require.Error(t, err)
				assert.Equal(t, tc.expectedErr, err.Error())
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, "alat-berat", resp.Slug)
				assert.Equal(t, catId, resp.Id)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCategoryService_EditAttribute(t *testing.T) {
	attrId := uuid.New()
	catId := uuid.New()

	existingAttr := &CategoryAttributes{
		BaseEntity: entity.BaseEntity{Id: attrId},
		CategoryId: catId,
		Name:       "Tahun",
		DataType:   "NUMBER",
	}

	existingCat := &Category{
		BaseEntity: entity.BaseEntity{Id: catId},
		Name:       "Mobil",
		Slug:       "mobil",
	}

	validReq := &UpdateCatAttrRequest{
		Id:         attrId,
		CategoryId: catId,
		Name:       "Tahun Pembuatan",
		DataType:   "NUMBER",
		Unit:       "Tahun",
		Options:    []string{},
		IsRequired: true,
		ShowOnCard: true,
		IsActive:   true,
		SortOrder:  1,
	}

	tests := []struct {
		name        string
		req         *UpdateCatAttrRequest
		setupMock   func(m *mockCategoryRepository)
		wantErr     bool
		expectedErr string
	}{
		{
			name: "gagal ketika attribute tidak ditemukan (nil, nil)",
			req:  validReq,
			setupMock: func(m *mockCategoryRepository) {
				// Menguji perbaikan bug operator &&: saat return (nil, nil), wajib langsung stop!
				m.On("FindAttributeById", mock.Anything, attrId).Return(nil, nil)
			},
			wantErr:     true,
			expectedErr: "attribute not found",
		},
		{
			name: "gagal ketika attribute ada tapi category_id tidak ditemukan",
			req:  validReq,
			setupMock: func(m *mockCategoryRepository) {
				m.On("FindAttributeById", mock.Anything, attrId).Return(existingAttr, nil)
				m.On("FindById", mock.Anything, catId).Return(nil, nil)
			},
			wantErr:     true,
			expectedErr: "category not found",
		},
		{
			name: "sukses update attribute saat attribute dan category valid",
			req:  validReq,
			setupMock: func(m *mockCategoryRepository) {
				m.On("FindAttributeById", mock.Anything, attrId).Return(existingAttr, nil)
				m.On("FindById", mock.Anything, catId).Return(existingCat, nil)
				m.On("UpdateAttribute", mock.Anything, mock.AnythingOfType("CategoryAttributes")).Return(&CategoryAttributes{
					BaseEntity: entity.BaseEntity{Id: attrId},
					CategoryId: catId,
					Name:       "Tahun Pembuatan",
					DataType:   "NUMBER",
					Unit:       "Tahun",
					Options:    []string{},
					IsRequired: true,
					ShowOnCard: true,
					IsActive:   true,
					SortOrder:  1,
				}, nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mockCategoryRepository)
			tc.setupMock(mockRepo)

			svc := NewCategoryService(mockRepo)
			resp, err := svc.EditAttribute(context.Background(), tc.req)

			if tc.wantErr {
				require.Error(t, err)
				assert.Equal(t, tc.expectedErr, err.Error())
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, "Tahun Pembuatan", resp.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

