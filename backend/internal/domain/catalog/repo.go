package catalog

import (
	"context"
	"errors"

	"chatsmith-mini/pkg/databases"
	"chatsmith-mini/pkg/failure"

	"gorm.io/gorm"
)

type Repo struct {
	db databases.PostgresDB
}

func newRepo(db databases.PostgresDB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) List(ctx context.Context) ([]Model, error) {
	var items []Model
	if err := r.db.WithContext(ctx).Where("enabled = ?", true).Order("sort_order, id").Find(&items).Error; err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return items, nil
}

func (r *Repo) Get(ctx context.Context, id string) (*Model, error) {
	var item Model
	err := r.db.WithContext(ctx).Where("enabled = ? AND id = ?", true, id).Take(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, failure.New(failure.CodeModelUnavailable)
	}
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return &item, nil
}
