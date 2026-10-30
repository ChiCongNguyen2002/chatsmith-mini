package user

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

func (r *Repo) Get(ctx context.Context, userID string) (*User, error) {
	var item User
	err := r.db.WithContext(ctx).Where("id = ?", userID).Take(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, failure.New(failure.CodeUnauthorized)
	}
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return &item, nil
}

func (r *Repo) FindByIdentity(ctx context.Context, provider, subject string) (*User, error) {
	return r.findJoined(ctx, "user_identities.provider = ? AND user_identities.subject = ?", provider, subject)
}

func (r *Repo) FindByVerifiedEmail(ctx context.Context, email string) (*User, error) {
	return r.findJoined(ctx, "user_identities.email = ? AND user_identities.email_verified = true", email)
}

func (r *Repo) findJoined(ctx context.Context, where string, args ...any) (*User, error) {
	var item User
	err := r.db.WithContext(ctx).Model(&User{}).Select("users.*").
		Joins("JOIN user_identities ON user_identities.user_id = users.id").Where(where, args...).Take(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return &item, nil
}

func (r *Repo) CreateWithIdentity(ctx context.Context, identity Identity) (*User, error) {
	var item User
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		item = User{Email: identity.Email, Name: identity.Name, Picture: identity.Picture}
		if err := tx.Create(&item).Error; err != nil {
			return err
		}

		row := identityFrom(item.ID, identity)

		return tx.Create(&row).Error
	})
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return &item, nil
}

func (r *Repo) AddIdentity(ctx context.Context, userID string, identity Identity) error {
	row := identityFrom(userID, identity)
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return failure.ErrWithTrace(err)
	}

	return nil
}

func (r *Repo) RefreshProfile(ctx context.Context, userID string, identity Identity) error {
	err := r.db.WithContext(ctx).Exec("UPDATE users SET email = ?, name = ?, picture = ? WHERE id = ?",
		identity.Email, identity.Name, identity.Picture, userID).Error
	if err != nil {
		return failure.ErrWithTrace(err)
	}

	return nil
}

func identityFrom(userID string, identity Identity) identityRow {
	return identityRow{
		Provider:      string(identity.Provider),
		Subject:       identity.Subject,
		UserID:        userID,
		Email:         identity.Email,
		EmailVerified: identity.EmailVerified,
	}
}
