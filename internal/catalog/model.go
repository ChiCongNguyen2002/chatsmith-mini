package catalog

import "chatsmith-mini/constants"

type Model struct {
	ID          string             `gorm:"column:id;primaryKey"`
	Provider    constants.Provider `gorm:"column:provider"`
	DisplayName string             `gorm:"column:display_name"`
	Description string             `gorm:"column:description"`
	InputMicro  int64              `gorm:"column:input_micro_usd_per_mtok"`
	OutputMicro int64              `gorm:"column:output_micro_usd_per_mtok"`
	Enabled     bool               `gorm:"column:enabled"`
	IsNew       bool               `gorm:"column:is_new"`
	SortOrder   int                `gorm:"column:sort_order"`
}

func (Model) TableName() string {
	return tableModels
}
