package conversation

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

func (r *Repo) Create(ctx context.Context, userID, title string) (*Conversation, error) {
	created := Conversation{UserID: userID, Title: title}
	err := r.db.WithContext(ctx).Create(&created).Error
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return &created, nil
}

func (r *Repo) Get(ctx context.Context, conversationID, userID string) (*Conversation, error) {
	var item Conversation
	err := r.db.WithContext(ctx).Where(whereOwnedConversation, conversationID, userID).Take(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, failure.New(failure.CodeNotFound)
	}
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return &item, nil
}

func (r *Repo) List(ctx context.Context, userID string, limit int, after *cursor) ([]Conversation, error) {
	query := r.db.WithContext(ctx).Where(whereUserConversations, userID).Order(orderNewestFirst).Limit(limit)
	if after != nil {
		query = query.Where(whereBeforeCursor, after.updatedAt, after.id)
	}

	items := make([]Conversation, 0, limit)
	if err := query.Find(&items).Error; err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return items, nil
}

func (r *Repo) Delete(ctx context.Context, conversationID, userID string) error {
	result := r.db.WithContext(ctx).Where(whereOwnedConversation, conversationID, userID).Delete(&Conversation{})
	if result.Error != nil {
		return failure.ErrWithTrace(result.Error)
	}
	if result.RowsAffected == 0 {
		return failure.New(failure.CodeNotFound)
	}

	return nil
}

func (r *Repo) UpdateTitle(ctx context.Context, conversationID, userID, title string) error {
	result := r.db.WithContext(ctx).Model(&Conversation{}).Where(whereOwnedConversation, conversationID, userID).Update("title", title)
	if result.Error != nil {
		return failure.ErrWithTrace(result.Error)
	}
	if result.RowsAffected == 0 {
		return failure.New(failure.CodeNotFound)
	}

	return nil
}

func (r *Repo) Messages(ctx context.Context, conversationID string, afterSeq int64, limit int) ([]Message, error) {
	items := make([]Message, 0, limit)
	err := r.db.WithContext(ctx).Where(whereMessagesAfterSeq, conversationID, afterSeq).Order(orderOldestFirst).Limit(limit).Find(&items).Error
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return items, nil
}
