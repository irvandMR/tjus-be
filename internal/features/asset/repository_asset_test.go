package asset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapSortOrder(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "sort termurah (price_asc)",
			input:    "price_asc",
			expected: "price ASC, created_at DESC",
		},
		{
			name:     "sort termahal (price_desc)",
			input:    "price_desc",
			expected: "price DESC, created_at DESC",
		},
		{
			name:     "sort terpopuler dengan huruf besar dan spasi",
			input:    "  POPULAR  ",
			expected: "view_count DESC NULLS LAST, created_at DESC",
		},
		{
			name:     "sort terlama (oldest)",
			input:    "oldest",
			expected: "created_at ASC",
		},
		{
			name:     "sort terbaru (newest)",
			input:    "newest",
			expected: "created_at DESC",
		},
		{
			name:     "kosong kembali ke default newest",
			input:    "",
			expected: "created_at DESC",
		},
		{
			name:     "percobaan SQL Injection harus aman balik ke default",
			input:    "price ASC; DROP TABLE tbl_assets;--",
			expected: "created_at DESC",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := mapSortOrder(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected string
	}{
		{
			name:     "judul standar huruf besar kecil",
			title:    "Mazda CX-3 2019",
			expected: "mazda-cx-3-2019",
		},
		{
			name:     "judul dengan spasi di awal dan akhir",
			title:    "  Rumah Mewah BSD  ",
			expected: "rumah-mewah-bsd",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := generateSlug(tc.title)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

