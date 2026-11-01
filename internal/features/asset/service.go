package asset

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
	"tjus-be/internal/features/category"
	"tjus-be/internal/pkg/entity"
	"tjus-be/internal/pkg/util"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const assetCacheTTL = 30 * time.Minute

type AssetService struct {
	assetRepo   AssetRepository
	catRepo     category.CategoryRepository
	redisClient *redis.Client
}

func NewAssetService(assetRepo AssetRepository, catRepo category.CategoryRepository, redisClient *redis.Client) *AssetService {
	return &AssetService{
		assetRepo:   assetRepo,
		catRepo:     catRepo,
		redisClient: redisClient,
	}
}

func assetSlugCacheKey(slug string) string {
	return "asset:slug:" + slug
}

func (s *AssetService) invalidateAssetCache(ctx context.Context, slug string) {
	if s.redisClient == nil || slug == "" {
		return
	}
	cacheKey := assetSlugCacheKey(slug)
	if err := util.DeleteCache(ctx, s.redisClient, cacheKey); err != nil {
		log.Printf("[REDIS] Gagal hapus cache %s: %v", cacheKey, err)
	} else {
		log.Printf("[CACHE INVALIDATED] %s dihapus dari Redis", cacheKey)
	}
}

func (s *AssetService) validateCategoryAndAttributes(ctx context.Context, categoryId uuid.UUID, attrs []AttributeRequest) error {
	cat, err := s.catRepo.FindById(ctx, categoryId)
	if err != nil {
		return err
	}
	if cat == nil {
		return errors.New("category not found")
	}

	if len(attrs) > 0 {
		validAttrs, err := s.catRepo.FindAttributesByCategoryId(ctx, categoryId)
		if err != nil {
			return err
		}

		validMap := make(map[uuid.UUID]bool, len(validAttrs))
		for _, va := range validAttrs {
			validMap[va.Id] = true
		}

		for _, reqAttr := range attrs {
			if !validMap[reqAttr.CategoryAttributeId] {
				if strings.TrimSpace(reqAttr.CustomeLabel) != "" {
					return fmt.Errorf("attribute '%s' not found for category '%s'", reqAttr.CustomeLabel, cat.Name)
				}
				return fmt.Errorf("attribute not found for category '%s'", cat.Name)
			}
		}
	}

	return nil
}

// CreateAsset menangani orkestrasi pembuatan aset, pemilik (owner), dan atribut sekaligus
func (s *AssetService) CreateAsset(ctx context.Context, req *AssetRequest) (*AssetResponse, error) {
	// 1. Validasi input dasar & relasi kategori/atribut (Fail-Fast sebelum buka transaksi)
	if strings.TrimSpace(req.Title) == "" {
		return nil, errors.New("title cannot be empty")
	}
	if err := s.validateCategoryAndAttributes(ctx, req.CategoryId, req.Attributes); err != nil {
		return nil, err
	}

	// 2. Generate Kode Unik dan Slug
	code := fmt.Sprintf("AST-%d", time.Now().Unix())
	baseSlug := generateSlug(req.Title)
	finalSlug := baseSlug

	// Pastikan slug unik, jika sudah ada tambahkan suffix timestamp
	existingAsset, err := s.assetRepo.FindBySlug(ctx, finalSlug)
	if err != nil {
		return nil, err
	}
	if existingAsset != nil {
		finalSlug = fmt.Sprintf("%s-%d", baseSlug, time.Now().Unix()%10000)
	}

	// 3. Buka transaksi database secara langsung dan eksplisit
	tx, err := s.assetRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // otomatis rollback jika error atau belum di-commit

	// 4. Simpan data induk Asset
	newAsset, err := s.assetRepo.InsertAssetTx(ctx, tx, &Asset{
		Code:            code,
		Slug:            finalSlug,
		CategoryId:      req.CategoryId,
		Title:           req.Title,
		Location:        req.Location,
		Description:     req.Description,
		Badge:           req.Badge,
		PriceLabel:      req.PriceLabel,
		Price:           req.Price,
		PriceNegotiable: req.PriceNegotiable,
		ViewCount:       0,
		Status:          "ACTIVE",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert asset: %w", err)
	}

	// 5. Simpan data Asset Owner
	owner := &AssetOwner{
		AssetId:       newAsset.Id,
		OwnerName:     req.Owner.OwnerName,
		OwnerContact:  req.Owner.OwnerContact,
		InternalNotes: req.Owner.InternalNotes,
	}
	if err := s.assetRepo.InsertAssetOwnerTx(ctx, tx, owner); err != nil {
		return nil, fmt.Errorf("failed to insert asset owner: %w", err)
	}

	// 6. Simpan data Atribut (jika ada) menggunakan Bulk Insert
	if len(req.Attributes) > 0 {
		var attributeEntities []AssetAttributes
		for _, attr := range req.Attributes {
			attributeEntities = append(attributeEntities, AssetAttributes{
				AssetId:             newAsset.Id,
				CategoryAttributeId: attr.CategoryAttributeId,
				CustomeLabel:        attr.CustomeLabel,
				ValueText:           attr.ValueText,
				ValueNumber:         attr.ValueNumber,
				SortOrder:           attr.SortOrder,
			})
		}

		if err := s.assetRepo.BulkInsertAssetAttributesTx(ctx, tx, attributeEntities); err != nil {
			return nil, fmt.Errorf("failed to insert asset attributes: %w", err)
		}
	}

	// 7. Commit seluruh perubahan jika semua query sukses!
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 8. Ambil data atribut lengkap (JOIN ke master kategori attributes untuk nama, tipe data, dan satuan unit)
	var attributeResponses []AssetAttributeResponse
	if len(req.Attributes) > 0 {
		savedAttrs, err := s.assetRepo.FindAttributesByAssetId(ctx, newAsset.Id)
		if err == nil {
			for _, attr := range savedAttrs {
				attributeResponses = append(attributeResponses, AssetAttributeResponse{
					CategoryAttributeId: attr.CategoryAttributeId,
					Name:                attr.Name,
					DataType:            attr.DataType,
					Unit:                attr.Unit,
					CustomeLabel:        attr.CustomeLabel,
					ValueText:           attr.ValueText,
					ValueNumber:         attr.ValueNumber,
					SortOrder:           attr.SortOrder,
				})
			}
		}
	}

	// 6. Susun dan kembalikan Response
	return &AssetResponse{
		Id:              newAsset.Id,
		Code:            newAsset.Code,
		Slug:            newAsset.Slug,
		CategoryId:      newAsset.CategoryId,
		Title:           newAsset.Title,
		Location:        newAsset.Location,
		Description:     newAsset.Description,
		Badge:           newAsset.Badge,
		PriceLabel:      newAsset.PriceLabel,
		Price:           newAsset.Price,
		PriceNegotiable: newAsset.PriceNegotiable,
		ViewCount:       newAsset.ViewCount,
		Status:          newAsset.Status,
		CreatedAt:       newAsset.CreatedAt,
		UpdatedAt:       newAsset.UpdatedAt,
		Owner: &AssetOwnerResponse{
			OwnerName:     owner.OwnerName,
			OwnerContact:  owner.OwnerContact,
			InternalNotes: owner.InternalNotes,
		},
		Attributes: attributeResponses,
		Media:      []AssetMediaResponse{}, // Media kosong saat awal pembuatan
	}, nil
}

// GetAssetById mengambil data lengkap aset berdasarkan ID (detail view untuk Admin - termasuk data owner)
func (s *AssetService) GetAssetById(ctx context.Context, id uuid.UUID) (*AssetResponse, error) {
	asset, err := s.assetRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, errors.New("asset not found")
	}

	return s.mapToAssetResponse(ctx, asset, true)
}

// GetAssetBySlug mengambil data lengkap aset berdasarkan Slug (halaman publik pengunjung - TANPA data owner)
func (s *AssetService) GetAssetBySlug(ctx context.Context, slug string) (*AssetResponse, error) {
	start := time.Now()
	cacheKey := assetSlugCacheKey(slug)

	// 1. Cek di RAM Redis terlebih dahulu (Cache-Aside)
	var cachedResp AssetResponse
	if util.GetCache(ctx, s.redisClient, cacheKey, &cachedResp) {
		log.Printf("[CACHE HIT - Redis] %s selesai dalam %v", cacheKey, time.Since(start))
		go s.assetRepo.IncrementViewCount(context.Background(), cachedResp.Id)
		cachedResp.ViewCount++
		return &cachedResp, nil
	}

	// 2. Cache Miss: Ambil dari PostgreSQL
	asset, err := s.assetRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if asset == nil || asset.Status != "ACTIVE" {
		return nil, errors.New("asset not found")
	}

	if err := s.assetRepo.IncrementViewCount(ctx, asset.Id); err != nil {
		return nil, err
	}
	asset.ViewCount++

	// includeOwner = false: data nomor HP dan catatan rahasia pemilik disembunyikan dari publik
	resp, err := s.mapToAssetResponse(ctx, asset, false)
	if err != nil {
		return nil, err
	}

	log.Printf("[CACHE MISS - Postgres] %s selesai dalam %v", cacheKey, time.Since(start))

	// 3. Simpan hasil ke Redis dengan TTL 30 menit
	_ = util.SetCache(ctx, s.redisClient, cacheKey, resp, assetCacheTTL)

	return resp, nil
}

// mapToAssetResponse helper untuk memetakan entity aset dengan atribut, media, dan opsi owner
func (s *AssetService) mapToAssetResponse(ctx context.Context, asset *Asset, includeOwner bool) (*AssetResponse, error) {
	var ownerResp *AssetOwnerResponse
	if includeOwner {
		owner, err := s.assetRepo.FindOwnerByAssetId(ctx, asset.Id)
		if err != nil {
			return nil, err
		}
		if owner != nil {
			ownerResp = &AssetOwnerResponse{
				OwnerName:     owner.OwnerName,
				OwnerContact:  owner.OwnerContact,
				InternalNotes: owner.InternalNotes,
			}
		}
	}

	// Ambil Atribut
	attributes, err := s.assetRepo.FindAttributesByAssetId(ctx, asset.Id)
	if err != nil {
		return nil, err
	}

	var attrResp []AssetAttributeResponse
	for _, attr := range attributes {
		attrResp = append(attrResp, AssetAttributeResponse{
			CategoryAttributeId: attr.CategoryAttributeId,
			Name:                attr.Name,
			DataType:            attr.DataType,
			Unit:                attr.Unit,
			CustomeLabel:        attr.CustomeLabel,
			ValueText:           attr.ValueText,
			ValueNumber:         attr.ValueNumber,
			SortOrder:           attr.SortOrder,
		})
	}

	// Ambil Media
	mediaList, err := s.assetRepo.FindMediaByAssetId(ctx, asset.Id)
	if err != nil {
		return nil, err
	}

	var mediaResp []AssetMediaResponse
	for _, m := range mediaList {
		mediaResp = append(mediaResp, AssetMediaResponse{
			Id:        m.Id,
			TypeMedia: m.TypeMedia,
			FilePath:  m.FilePath,
			SortOrder: m.SortOrder,
		})
	}

	return &AssetResponse{
		Id:              asset.Id,
		Code:            asset.Code,
		Slug:            asset.Slug,
		CategoryId:      asset.CategoryId,
		Title:           asset.Title,
		Location:        asset.Location,
		Description:     asset.Description,
		Badge:           asset.Badge,
		PriceLabel:      asset.PriceLabel,
		Price:           asset.Price,
		PriceNegotiable: asset.PriceNegotiable,
		ViewCount:       asset.ViewCount,
		Status:          asset.Status,
		CreatedAt:       asset.CreatedAt,
		UpdatedAt:       asset.UpdatedAt,
		Owner:           ownerResp,
		Attributes:      attrResp,
		Media:           mediaResp,
	}, nil
}

func generateSlug(text string) string {
	lower := strings.ToLower(strings.TrimSpace(text))
	return strings.ReplaceAll(lower, " ", "-")
}

// GetAllAssets mengambil daftar aset aktif dengan dukungan pencarian, filtering, dan pagination
func (s *AssetService) GetAllAssets(ctx context.Context, filter AssetFilter) ([]AssetResponse, util.Pagination, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	} else if filter.Limit > 100 {
		filter.Limit = 100
	}

	filter.Offset = (filter.Page - 1) * filter.Limit

	// Public security guard: paksa status selalu ACTIVE untuk katalog pengunjung umum
	active := "ACTIVE"
	filter.Status = &active

	assets, total, err := s.assetRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, util.Pagination{}, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(total / int64(filter.Limit))
		if total%int64(filter.Limit) != 0 {
			totalPages++
		}
	}

	pagination := util.Pagination{
		CurrentPage: filter.Page,
		Limit:       filter.Limit,
		TotalItems:  total,
		TotalPages:  totalPages,
	}

	assetResponses := make([]AssetResponse, 0, len(assets))
	for _, a := range assets {
		// Ambil media (gambar/thumbnail) untuk preview aset pada list catalog
		mediaList, err := s.assetRepo.FindMediaByAssetId(ctx, a.Id)
		if err != nil {
			return nil, util.Pagination{}, err
		}

		mediaResp := make([]AssetMediaResponse, 0, len(mediaList))
		for _, m := range mediaList {
			mediaResp = append(mediaResp, AssetMediaResponse{
				Id:        m.Id,
				TypeMedia: m.TypeMedia,
				FilePath:  m.FilePath,
				SortOrder: m.SortOrder,
			})
		}

		assetResponses = append(assetResponses, AssetResponse{
			Id:              a.Id,
			Code:            a.Code,
			Slug:            a.Slug,
			CategoryId:      a.CategoryId,
			Title:           a.Title,
			Location:        a.Location,
			Description:     a.Description,
			Badge:           a.Badge,
			PriceLabel:      a.PriceLabel,
			Price:           a.Price,
			PriceNegotiable: a.PriceNegotiable,
			ViewCount:       a.ViewCount,
			Status:          a.Status,
			CreatedAt:       a.CreatedAt,
			UpdatedAt:       a.UpdatedAt,
			Media:           mediaResp,
		})
	}

	return assetResponses, pagination, nil
}

// GetAdminAssets mengambil daftar aset untuk dashboard Admin (bisa filter semua status dan menyertakan data owner)
func (s *AssetService) GetAdminAssets(ctx context.Context, filter AssetFilter) ([]AssetResponse, util.Pagination, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	} else if filter.Limit > 100 {
		filter.Limit = 100
	}

	filter.Offset = (filter.Page - 1) * filter.Limit

	// Catatan: untuk admin, filter.Status TIDAK dipaksa ke ACTIVE. Admin bebas memfilter status apapun.
	assets, total, err := s.assetRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, util.Pagination{}, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(total / int64(filter.Limit))
		if total%int64(filter.Limit) != 0 {
			totalPages++
		}
	}

	pagination := util.Pagination{
		CurrentPage: filter.Page,
		Limit:       filter.Limit,
		TotalItems:  total,
		TotalPages:  totalPages,
	}

	assetResponses := make([]AssetResponse, 0, len(assets))
	for _, a := range assets {
		// Ambil data owner untuk dashboard admin
		var ownerResp *AssetOwnerResponse
		owner, err := s.assetRepo.FindOwnerByAssetId(ctx, a.Id)
		if err == nil && owner != nil {
			ownerResp = &AssetOwnerResponse{
				OwnerName:     owner.OwnerName,
				OwnerContact:  owner.OwnerContact,
				InternalNotes: owner.InternalNotes,
			}
		}

		// Ambil media thumbnail/preview
		mediaList, _ := s.assetRepo.FindMediaByAssetId(ctx, a.Id)
		mediaResp := make([]AssetMediaResponse, 0, len(mediaList))
		for _, m := range mediaList {
			mediaResp = append(mediaResp, AssetMediaResponse{
				Id:        m.Id,
				TypeMedia: m.TypeMedia,
				FilePath:  m.FilePath,
				SortOrder: m.SortOrder,
			})
		}

		assetResponses = append(assetResponses, AssetResponse{
			Id:              a.Id,
			Code:            a.Code,
			Slug:            a.Slug,
			CategoryId:      a.CategoryId,
			Title:           a.Title,
			Location:        a.Location,
			Description:     a.Description,
			Badge:           a.Badge,
			PriceLabel:      a.PriceLabel,
			Price:           a.Price,
			PriceNegotiable: a.PriceNegotiable,
			ViewCount:       a.ViewCount,
			Status:          a.Status,
			CreatedAt:       a.CreatedAt,
			UpdatedAt:       a.UpdatedAt,
			Owner:           ownerResp,
			Media:           mediaResp,
		})
	}

	return assetResponses, pagination, nil
}

// UploadMedia menangani validasi ketat (ekstensi, ukuran, MIME sniffing), penyimpanan file ke storage lokal, dan pencatatan ke database tbl_asset_media.
// Mendukung upload satu file maupun banyak file sekaligus (batch upload).
func (s *AssetService) UploadMedia(ctx context.Context, assetId uuid.UUID, files []*multipart.FileHeader) ([]AssetMediaResponse, error) {
	// 1. Cek apakah aset-nya ada di database
	asset, err := s.assetRepo.FindById(ctx, assetId)
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, errors.New("asset not found")
	}

	// 2. Validasi seluruh file secara dinamis menggunakan fileutil (All-or-Nothing sebelum disimpan ke disk)
	err = util.ValidateUploadedFiles(files, util.FileValidationConfig{
		MaxFiles:     10,
		MaxFileSize:  10 * 1024 * 1024, // 10 MB per foto
		AllowedExts:  map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true},
		AllowedMimes: map[string]bool{"image/jpeg": true, "image/png": true, "image/webp": true},
	})
	if err != nil {
		return nil, err
	}

	// 3. Ambil media yang sudah ada untuk menentukan urutan (sort_order) lanjutan
	existingMedia, err := s.assetRepo.FindMediaByAssetId(ctx, assetId)
	if err != nil {
		return nil, err
	}
	currentSortOrder := len(existingMedia)

	uploadDir := "uploads/assets"
	var responses []AssetMediaResponse

	// 4. Eksekusi penyimpanan file ke disk & database
	for i, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		fileName := fmt.Sprintf("%s-%d-%d%s", assetId.String()[:8], time.Now().UnixNano(), i, ext)
		fullLocalPath := filepath.Join(uploadDir, fileName)

		if err := util.SaveUploadedFile(file, fullLocalPath); err != nil {
			return nil, err
		}

		// Simpan record ke database tbl_asset_media
		publicUrl := "/uploads/assets/" + fileName
		sortOrder := currentSortOrder + i + 1

		mediaRecord, err := s.assetRepo.InsertAssetMedia(ctx, &AssetMedia{
			AssetId:   assetId,
			TypeMedia: "IMAGE",
			FilePath:  publicUrl,
			SortOrder: sortOrder,
		})
		if err != nil {
			_ = os.Remove(fullLocalPath)
			return nil, fmt.Errorf("failed to insert media record for '%s': %w", file.Filename, err)
		}

		responses = append(responses, AssetMediaResponse{
			Id:        mediaRecord.Id,
			TypeMedia: mediaRecord.TypeMedia,
			FilePath:  mediaRecord.FilePath,
			SortOrder: mediaRecord.SortOrder,
		})
	}

	s.invalidateAssetCache(ctx, asset.Slug)
	return responses, nil
}

// UpdateAsset menangani pembaruan aset secara modular berdasarkan boolean flag
func (s *AssetService) UpdateAsset(ctx context.Context, req *UpdateAssetRequest) (*AssetResponse, error) {
	// 1. Cek apakah ada minimal satu flag update yang aktif
	if !req.GeneralUpdate && !req.OwnerUpdate && !req.AttributeUpdate && !req.StatusUpdate {
		return nil, errors.New("at least one update flag must be true (general_update, owner_update, attribute_update, status_update)")
	}

	// 2. Cek apakah aset-nya ada di database
	existingAsset, err := s.assetRepo.FindById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if existingAsset == nil {
		return nil, errors.New("asset not found")
	}

	// 2b. Jika mengupdate atribut, validasi dulu atributnya cocok dengan kategori aset sebelum buka transaksi
	if req.AttributeUpdate && len(req.Attributes) > 0 {
		if err := s.validateCategoryAndAttributes(ctx, existingAsset.CategoryId, req.Attributes); err != nil {
			return nil, err
		}
	}

	// 3. Buka transaksi database
	tx, err := s.assetRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// 4. Update Status jika status_update == true
	if req.StatusUpdate {
		status := strings.TrimSpace(req.Status)
		if status == "" {
			return nil, errors.New("status cannot be empty when status_update is true")
		}
		if err := s.assetRepo.UpdateAssetStatusTx(ctx, tx, req.Id, status); err != nil {
			return nil, fmt.Errorf("failed to update asset status: %w", err)
		}
	}

	// 5. Update Informasi Umum jika general_update == true
	if req.GeneralUpdate {
		if strings.TrimSpace(req.Title) == "" {
			return nil, errors.New("title cannot be empty when general_update is true")
		}

		slug := existingAsset.Slug
		// Jika judul berubah, update slug
		if req.Title != existingAsset.Title {
			baseSlug := generateSlug(req.Title)
			slug = baseSlug
			existingSlug, err := s.assetRepo.FindBySlug(ctx, slug)
			if err == nil && existingSlug != nil && existingSlug.Id != req.Id {
				slug = fmt.Sprintf("%s-%d", baseSlug, time.Now().Unix()%10000)
			}
		}

		assetToUpdate := &Asset{
			BaseEntity: entity.BaseEntity{
				Id: req.Id,
			},
			CategoryId:      existingAsset.CategoryId,
			Title:           req.Title,
			Slug:            slug,
			Location:        req.Location,
			Description:     req.Description,
			Badge:           req.Badge,
			PriceLabel:      req.PriceLabel,
			Price:           req.Price,
			PriceNegotiable: req.PriceNegotiable,
		}

		if err := s.assetRepo.UpdateAssetTx(ctx, tx, assetToUpdate); err != nil {
			return nil, fmt.Errorf("failed to update asset info: %w", err)
		}
	}

	// 6. Update Owner Info jika owner_update == true
	if req.OwnerUpdate {
		ownerToUpdate := &AssetOwner{
			AssetId:       req.Id,
			OwnerName:     req.Owner.OwnerName,
			OwnerContact:  req.Owner.OwnerContact,
			InternalNotes: req.Owner.InternalNotes,
		}
		if err := s.assetRepo.UpdateAssetOwnerTx(ctx, tx, ownerToUpdate); err != nil {
			return nil, fmt.Errorf("failed to update asset owner: %w", err)
		}
	}

	// 7. Update Atribut jika attribute_update == true (clean replace)
	if req.AttributeUpdate {
		if err := s.assetRepo.DeleteAttributesByAssetIdTx(ctx, tx, req.Id); err != nil {
			return nil, fmt.Errorf("failed to clear old attributes: %w", err)
		}

		if len(req.Attributes) > 0 {
			var attributeEntities []AssetAttributes
			for _, attr := range req.Attributes {
				attributeEntities = append(attributeEntities, AssetAttributes{
					AssetId:             req.Id,
					CategoryAttributeId: attr.CategoryAttributeId,
					CustomeLabel:        attr.CustomeLabel,
					ValueText:           attr.ValueText,
					ValueNumber:         attr.ValueNumber,
					SortOrder:           attr.SortOrder,
				})
			}

			if err := s.assetRepo.BulkInsertAssetAttributesTx(ctx, tx, attributeEntities); err != nil {
				return nil, fmt.Errorf("failed to insert new attributes: %w", err)
			}
		}
	}

	// 8. Commit transaksi
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 9. Hapus cache lama di Redis agar pengunjung langsung mendapat data terbaru
	s.invalidateAssetCache(ctx, existingAsset.Slug)

	// 10. Kembalikan data aset terupdate secara lengkap
	return s.GetAssetById(ctx, req.Id)
}

// DeleteMedia menghapus data media dari database sekaligus menghapus file fisiknya dari storage lokal
func (s *AssetService) DeleteMedia(ctx context.Context, mediaId uuid.UUID) error {
	// 1. Cek apakah media ada di database
	media, err := s.assetRepo.FindMediaById(ctx, mediaId)
	if err != nil {
		return err
	}
	if media == nil {
		return errors.New("media not found")
	}

	// 2. Hapus file fisik dari storage lokal jika ada
	filePath := strings.TrimPrefix(media.FilePath, "/")
	_ = os.Remove(filePath)

	// 3. Hapus record dari database
	if err := s.assetRepo.DeleteMediaById(ctx, mediaId); err != nil {
		return fmt.Errorf("failed to delete media record: %w", err)
	}

	// 4. Hapus cache aset terkait di Redis
	if asset, err := s.assetRepo.FindById(ctx, media.AssetId); err == nil && asset != nil {
		s.invalidateAssetCache(ctx, asset.Slug)
	}

	return nil
}

// DeleteAsset melakukan soft delete aset dengan mengubah status menjadi DELETED
func (s *AssetService) DeleteAsset(ctx context.Context, id uuid.UUID) error {
	// 1. Cek apakah aset ada di database
	asset, err := s.assetRepo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if asset == nil || asset.Status == "DELETED" {
		return errors.New("asset not found")
	}

	// 2. Eksekusi soft delete di repository
	if err := s.assetRepo.DeleteAsset(ctx, id); err != nil {
		return err
	}

	// 3. Hapus cache aset dari Redis agar tidak lagi bisa diakses publik
	s.invalidateAssetCache(ctx, asset.Slug)
	return nil
}
