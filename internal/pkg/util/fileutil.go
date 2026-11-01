package util

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FileValidationConfig struct {
	MaxFiles     int
	MaxFileSize  int64
	AllowedExts  map[string]bool
	AllowedMimes map[string]bool
}

// ValidateUploadedFiles memeriksa seluruh file secara dinamis (All-or-Nothing) sebelum disimpan ke storage.
func ValidateUploadedFiles(files []*multipart.FileHeader, cfg FileValidationConfig) error {
	if cfg.MaxFiles > 0 && len(files) > cfg.MaxFiles {
		return fmt.Errorf("maximum %d files allowed per upload", cfg.MaxFiles)
	}

	for _, file := range files {
		if file.Size == 0 {
			return fmt.Errorf("file '%s' is empty (0 bytes)", file.Filename)
		}

		if cfg.MaxFileSize > 0 && file.Size > cfg.MaxFileSize {
			maxMB := cfg.MaxFileSize / (1024 * 1024)
			if maxMB > 0 {
				return fmt.Errorf("file '%s' exceeds maximum limit of %dMB", file.Filename, maxMB)
			}
			return fmt.Errorf("file '%s' exceeds maximum size limit", file.Filename)
		}

		if len(cfg.AllowedExts) > 0 {
			ext := strings.ToLower(filepath.Ext(file.Filename))
			if !cfg.AllowedExts[ext] {
				return fmt.Errorf("file '%s' has unsupported file extension '%s'", file.Filename, ext)
			}
		}

		if len(cfg.AllowedMimes) > 0 {
			detectedMime, err := detectFileMime(file)
			if err != nil {
				return err
			}
			if !cfg.AllowedMimes[detectedMime] {
				return fmt.Errorf("file '%s' content type is not allowed (detected: %s)", file.Filename, detectedMime)
			}
		}
	}

	return nil
}

func detectFileMime(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file '%s': %w", file.Filename, err)
	}
	defer src.Close()

	buf := make([]byte, 512)
	n, err := src.Read(buf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to inspect file content for '%s': %w", file.Filename, err)
	}

	return http.DetectContentType(buf[:n]), nil
}

// SaveUploadedFile menyimpan file multipart ke path tujuan secara aman dan menghapus file parsial jika gagal.
func SaveUploadedFile(file *multipart.FileHeader, dstPath string) error {
	if err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file '%s': %w", file.Filename, err)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to save file '%s' to disk: %w", file.Filename, err)
	}

	_, copyErr := io.Copy(dst, src)
	closeErr := dst.Close()

	if copyErr != nil {
		_ = os.Remove(dstPath)
		return fmt.Errorf("failed to write content for file '%s': %w", file.Filename, copyErr)
	}
	if closeErr != nil {
		_ = os.Remove(dstPath)
		return fmt.Errorf("failed to finalize file '%s': %w", file.Filename, closeErr)
	}

	return nil
}

