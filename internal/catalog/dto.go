package catalog

import "chatsmith-mini/constants"

type ModelView struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsNew       bool   `json:"is_new"`
	IsDefault   bool   `json:"is_default"`
	Available   bool   `json:"available"`
	InputMicro  int64  `json:"input_micro_usd_per_mtok"`
	OutputMicro int64  `json:"output_micro_usd_per_mtok"`
}

type ProviderView struct {
	ID        constants.Provider `json:"id"`
	Name      string             `json:"name"`
	Available bool               `json:"available"`
	Models    []ModelView        `json:"models"`
}

type Catalog struct {
	DefaultModel string         `json:"default_model"`
	Providers    []ProviderView `json:"providers"`
}
