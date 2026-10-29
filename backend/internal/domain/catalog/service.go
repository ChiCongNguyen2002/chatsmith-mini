package catalog

import (
	"context"
	"fmt"

	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/failure"
)

type Service struct {
	repo         *Repo
	available    func(constants.Provider) bool
	defaultModel string
}

func newService(repo *Repo, available func(constants.Provider) bool, defaultModel string) *Service {
	return &Service{repo: repo, available: available, defaultModel: defaultModel}
}

func (s *Service) List(ctx context.Context) (*Catalog, error) {
	models, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	byProvider := make(map[constants.Provider][]ModelView, len(providerOrder))
	for _, m := range models {
		byProvider[m.Provider] = append(byProvider[m.Provider], ModelView{
			ID:          m.ID,
			Name:        m.DisplayName,
			Description: m.Description,
			IsNew:       m.IsNew,
			IsDefault:   m.ID == s.defaultModel,
			Available:   s.available(m.Provider),
			InputMicro:  m.InputMicro,
			OutputMicro: m.OutputMicro,
		})
	}

	result := &Catalog{DefaultModel: s.defaultModel, Providers: make([]ProviderView, 0, len(providerOrder))}
	for _, provider := range providerOrder {
		if views := byProvider[provider]; len(views) > 0 {
			result.Providers = append(result.Providers, ProviderView{
				ID:        provider,
				Name:      providerNames[provider],
				Available: s.available(provider),
				Models:    views,
			})
		}
	}

	return result, nil
}

func (s *Service) Resolve(ctx context.Context, id string) (*Model, error) {
	if id == "" {
		id = s.defaultModel
	}

	model, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !s.available(model.Provider) {
		return nil, failure.New(failure.CodeModelUnavailable)
	}

	return model, nil
}

func (s *Service) CheckDefault(ctx context.Context) error {
	if _, err := s.repo.Get(ctx, s.defaultModel); err != nil {
		return fmt.Errorf("default model %q is not an enabled model: %w", s.defaultModel, err)
	}

	return nil
}
