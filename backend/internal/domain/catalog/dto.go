package catalog

import "chatsmith-mini/constants"

type ModelResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsNew       bool   `json:"is_new"`
	IsDefault   bool   `json:"is_default"`
	Available   bool   `json:"available"`
	InputMicro  int64  `json:"input_micro_usd_per_mtok"`
	OutputMicro int64  `json:"output_micro_usd_per_mtok"`
}

type ProviderResponse struct {
	ID        constants.Provider `json:"id"`
	Name      string             `json:"name"`
	Available bool               `json:"available"`
	Models    []ModelResponse    `json:"models"`
}

type CatalogResponse struct {
	DefaultModel string             `json:"default_model"`
	Providers    []ProviderResponse `json:"providers"`
}

func ToResponse(item Catalog) CatalogResponse {
	providers := make([]ProviderResponse, 0, len(item.Providers))
	for _, provider := range item.Providers {
		models := make([]ModelResponse, 0, len(provider.Models))
		for _, model := range provider.Models {
			models = append(models, ModelResponse{
				ID:          model.ID,
				Name:        model.Name,
				Description: model.Description,
				IsNew:       model.IsNew,
				IsDefault:   model.IsDefault,
				Available:   model.Available,
				InputMicro:  model.InputMicro,
				OutputMicro: model.OutputMicro,
			})
		}
		providers = append(providers, ProviderResponse{
			ID:        provider.ID,
			Name:      provider.Name,
			Available: provider.Available,
			Models:    models,
		})
	}

	return CatalogResponse{DefaultModel: item.DefaultModel, Providers: providers}
}
