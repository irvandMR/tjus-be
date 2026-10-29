package category

import (
	"tjus-be/internal/pkg/entity"

	"github.com/google/uuid"
)

type Category struct {
	entity.BaseEntity
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	SortOrder int    `json:"sort_order"`
	IsActive  bool   `json:"is_active"`
}

type CategoryAttributes struct {
	entity.BaseEntity
	CategoryId uuid.UUID `json:"category_id"`
	Name       string    `json:"name"`
	DataType   string    `json:"data_type"`
	Unit       string    `json:"unit"`
	Options    []string  `json:"options"`
	IsRequired bool      `json:"is_required"`
	ShowOnCard bool      `json:"show_on_card"`
	IsActive   bool      `json:"is_active"`
	SortOrder  int       `json:"sort_order"`
}

// Request

type CreatedCatRequest struct {
	Name      string `json:"name" binding:"required"`
	SortOrder int    `json:"sort_order" binding:"required"`
	IsActive  bool   `json:"is_active"`
}

type UpdateCatRequest struct {
	CreatedCatRequest
	Id uuid.UUID `json:"id"`
}

type AttributeItemRequest struct {
	Name       string   `json:"name" binding:"required"`
	DataType   string   `json:"data_type" binding:"required"`
	Unit       string   `json:"unit"`
	Options    []string `json:"options"`
	IsRequired bool     `json:"is_required"`
	ShowOnCard bool     `json:"show_on_card"`
	IsActive   bool     `json:"is_active"`
	SortOrder  int      `json:"sort_order"`
}

type CreatedCatAttrRequest struct {
	CategoryId uuid.UUID              `json:"category_id" binding:"required"`
	Attributes []AttributeItemRequest `json:"attributes" binding:"required,min=1"`
}

type UpdateCatAttrRequest struct {
	Id         uuid.UUID `json:"id" binding:"required"`
	CategoryId uuid.UUID `json:"category_id" binding:"required"`
	Name       string    `json:"name" binding:"required"`
	DataType   string    `json:"data_type" binding:"required"`
	Unit       string    `json:"unit"`
	Options    []string  `json:"options"`
	IsRequired bool      `json:"is_required"`
	ShowOnCard bool      `json:"show_on_card"`
	IsActive   bool      `json:"is_active"`
	SortOrder  int       `json:"sort_order" binding:"required"`
}

// Response

type CategoryResponse struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	SortOrder int       `json:"sort_order"`
	IsActive  bool      `json:"is_active"`
}

type CategoryAttributeResponse struct {
	Id         uuid.UUID `json:"id"`
	CategoryId uuid.UUID `json:"category_id"`
	Name       string    `json:"name"`
	DataType   string    `json:"data_type"`
	Unit       string    `json:"unit"`
	Options    []string  `json:"options"`
	IsRequired bool      `json:"is_required"`
	ShowOnCard bool      `json:"show_on_card"`
	IsActive   bool      `json:"is_active"`
	SortOrder  int       `json:"sort_order"`
}
