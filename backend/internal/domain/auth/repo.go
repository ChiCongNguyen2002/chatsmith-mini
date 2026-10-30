package auth

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

func (r *Repo) Create(ctx context.Context, session Session) error {
	if err := r.db.WithContext(ctx).Create(&session).Error; err != nil {
		return failure.ErrWithTrace(err)
	}

	return nil
}

func (r *Repo) Find(ctx context.Context, tokenHash string) (*Session, error) {
	var item Session
	err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).Take(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return &item, nil
}

func (r *Repo) Delete(ctx context.Context, tokenHash string) error {
	if err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).Delete(&Session{}).Error; err != nil {
		return failure.ErrWithTrace(err)
	}

	return nil
}
