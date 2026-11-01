package asset

import (
	"time"
	"tjus-be/internal/pkg/entity"

	"github.com/google/uuid"
)

// Entity

type Asset struct {
	entity.BaseEntity
	Code            string
	Slug            string
	CategoryId      uuid.UUID
	Title           string
	Location        string
	Description     string
	Badge           string
	PriceLabel      string
	Price           int64
	PriceNegotiable bool
	ViewCount       int32
	Status          string
}

type AssetAttributes struct {
	Id                  uuid.UUID
	AssetId             uuid.UUID
	CategoryAttributeId uuid.UUID
	Name                *string
	DataType            *string
	Unit                *string
	CustomeLabel        string
	ValueText           string
	ValueNumber         int64
	SortOrder           int
}

type AssetMedia struct {
	Id        uuid.UUID
	AssetId   uuid.UUID
	TypeMedia string
	FilePath  string
	SortOrder int
	CreatedAt time.Time
}

type AssetOwner struct {
	AssetId       uuid.UUID
	OwnerName     string
	OwnerContact  string
	InternalNotes string
	UpdatedAt     time.Time
}

// Request

type AssetRequest struct {
	CategoryId      uuid.UUID          `json:"category_id" binding:"required"`
	Title           string             `json:"title" binding:"required"`
	Location        string             `json:"location" binding:"required"`
	Badge           string             `json:"badge"`
	PriceLabel      string             `json:"price_label"`
	Price           int64              `json:"price" binding:"required"`
	PriceNegotiable bool               `json:"price_negotiable"`
	Owner           AssetOwnerRequest  `json:"owner"`
	Attributes      []AttributeRequest `json:"attributes"`
	Description     string             `json:"description"`
}

type AttributeRequest struct {
	CategoryAttributeId uuid.UUID `json:"attribute_id"`
	CustomeLabel        string    `json:"custome_label"`
	ValueText           string    `json:"value_text"`
	ValueNumber         int64     `json:"value_number"`
	SortOrder           int       `json:"sort_order"`
}

type AssetOwnerRequest struct {
	OwnerName     string `json:"owner_name"`
	OwnerContact  string `json:"owner_contact"`
	InternalNotes string `json:"internal_notes"`
}

type AssetMediaRequest struct {
	AssetId   uuid.UUID `json:"asset_id"`
	TypeMedia string    `json:"type_media"`
	FilePath  string    `json:"file_path"`
	SortOrder int       `json:"sort_order"`
}

type UpdateAssetRequest struct {
	Id              uuid.UUID          `json:"id" binding:"required"`
	GeneralUpdate   bool               `json:"general_update"`
	OwnerUpdate     bool               `json:"owner_update"`
	AttributeUpdate bool               `json:"attribute_update"`
	StatusUpdate    bool               `json:"status_update"`

	// Fields for GeneralUpdate
	Title           string             `json:"title"`
	Location        string             `json:"location"`
	Badge           string             `json:"badge"`
	PriceLabel      string             `json:"price_label"`
	Price           int64              `json:"price"`
	PriceNegotiable bool               `json:"price_negotiable"`
	Description     string             `json:"description"`

	// Field for StatusUpdate
	Status          string             `json:"status"`

	// Field for OwnerUpdate
	Owner           AssetOwnerRequest  `json:"owner"`

	// Field for AttributeUpdate
	Attributes      []AttributeRequest `json:"attributes"`
}

type AssetFilter struct {
	CategoryId *uuid.UUID `form:"category_id"`
	Search     string     `form:"q"`
	MinPrice   *int64     `form:"min_price"`
	MaxPrice   *int64     `form:"max_price"`
	Status     *string    `form:"status"`
	Sort       string     `form:"sort"`
	Page       int        `form:"page"`
	Limit      int        `form:"limit"`
	Offset     int        `form:"-"`
}

// Response

type AssetOwnerResponse struct {
	OwnerName     string `json:"owner_name"`
	OwnerContact  string `json:"owner_contact"`
	InternalNotes string `json:"internal_notes"`
}

type AssetAttributeResponse struct {
	CategoryAttributeId uuid.UUID `json:"category_attribute_id"`
	Name                *string   `json:"name,omitempty"`
	DataType            *string   `json:"data_type,omitempty"`
	Unit                *string   `json:"unit,omitempty"`
	CustomeLabel        string    `json:"custome_label"`
	ValueText           string    `json:"value_text"`
	ValueNumber         int64     `json:"value_number"`
	SortOrder           int       `json:"sort_order"`
}

type AssetMediaResponse struct {
	Id        uuid.UUID `json:"id"`
	TypeMedia string    `json:"type_media"`
	FilePath  string    `json:"file_path"`
	SortOrder int       `json:"sort_order"`
}

type AssetResponse struct {
	Id              uuid.UUID                `json:"id"`
	Code            string                   `json:"code"`
	Slug            string                   `json:"slug"`
	CategoryId      uuid.UUID                `json:"category_id"`
	Title           string                   `json:"title"`
	Location        string                   `json:"location"`
	Description     string                   `json:"description"`
	Badge           string                   `json:"badge"`
	PriceLabel      string                   `json:"price_label"`
	Price           int64                    `json:"price"`
	PriceNegotiable bool                     `json:"price_negotiable"`
	ViewCount       int32                    `json:"view_count"`
	Status          string                   `json:"status"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       *time.Time               `json:"updated_at,omitempty"`
	Owner           *AssetOwnerResponse      `json:"owner,omitempty"`
	Attributes      []AssetAttributeResponse `json:"attributes,omitempty"`
	Media           []AssetMediaResponse     `json:"media,omitempty"`
}
