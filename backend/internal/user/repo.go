package user

import (
	"context"
	"errors"

	"chatsmith-mini/pkg/databases"
	"chatsmith-mini/pkg/failure"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repo struct {
	db databases.PostgresDB
}

func newRepo(db databases.PostgresDB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Upsert(ctx context.Context, profile Profile) (*User, error) {
	item := User{GoogleSub: profile.GoogleSub, Email: profile.Email, Name: profile.Name, Picture: profile.Picture}
	err := r.db.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: colGoogleSub}},
			DoUpdates: clause.AssignmentColumns([]string{colEmail, colName, colPicture}),
		},
		clause.Returning{},
	).Create(&item).Error
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return &item, nil
}

func (r *Repo) Get(ctx context.Context, userID string) (*User, error) {
	var item User
	err := r.db.WithContext(ctx).Where(whereUserID, userID).Take(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, failure.New(failure.CodeUnauthorized)
	}
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return &item, nil
}
