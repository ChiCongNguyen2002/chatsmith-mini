package catalog

import "chatsmith-mini/constants"

type Model struct {
	Code          string             `gorm:"column:code"`
	ProviderModel string             `gorm:"column:provider_model"`
	Provider      constants.Provider `gorm:"column:provider"`
	DisplayName   string             `gorm:"column:display_name"`
	Description   string             `gorm:"column:description"`
	InputMicro    int64              `gorm:"column:input_micro_usd_per_mtok"`
	OutputMicro   int64              `gorm:"column:output_micro_usd_per_mtok"`
	IsActive      bool               `gorm:"column:is_active"`
	IsNew         bool               `gorm:"column:is_new"`
	OrderIndex    int                `gorm:"column:order_index"`
}

func (Model) TableName() string {
	return "catalog_models"
}

type ModelView struct {
	ID          string
	Name        string
	Description string
	IsNew       bool
	IsDefault   bool
	Available   bool
	InputMicro  int64
	OutputMicro int64
}

type ProviderView struct {
	ID        constants.Provider
	Name      string
	Available bool
	Models    []ModelView
}

type Catalog struct {
	DefaultModel string
	Providers    []ProviderView
}
