package chat

import (
	"context"
	"errors"
	"slices"
	"time"

	"chatsmith-mini/internal/conversation"
	"chatsmith-mini/pkg/databases"
	"chatsmith-mini/pkg/failure"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Repo struct {
	db databases.PostgresDB
}

func newRepo(db databases.PostgresDB) *Repo {
	return &Repo{db: db}
}

type finishInput struct {
	conversationID string
	messageID      string
	content        string
	status         conversation.MessageStatus
	usage          usageData
	costMicroUSD   int64
	firstTokenMS   *int
	latencyMS      int
}

func (r *Repo) Begin(ctx context.Context, conversationID, content, title, model string, staleBefore time.Time) (conversation.Message, conversation.Message, error) {
	user := conversation.Message{ConversationID: conversationID, Role: conversation.RoleUser, Content: content, Status: conversation.StatusCompleted}
	assistant := conversation.Message{ConversationID: conversationID, Role: conversation.RoleAssistant, Status: conversation.StatusStreaming, Model: &model}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := failStale(tx, conversationID, staleBefore); err != nil {
			return err
		}
		maxSeq, err := lastSeq(tx, conversationID)
		if err != nil {
			return err
		}
		user.Seq, assistant.Seq = maxSeq+1, maxSeq+2

		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if err := tx.Create(&assistant).Error; err != nil {
			return err
		}

		return tx.Exec(touchConversationSQL, title, conversationID).Error
	})
	if isUniqueViolation(err) {
		return user, assistant, failure.New(failure.CodeConversationBusy)
	}
	if err != nil {
		return user, assistant, failure.ErrWithTrace(err)
	}

	return user, assistant, nil
}

func (r *Repo) BeginContinue(ctx context.Context, truncated conversation.Message, staleBefore time.Time) (conversation.Message, error) {
	var affected int64
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := failStale(tx, truncated.ConversationID, staleBefore); err != nil {
			return err
		}
		res := tx.Exec(updateStatusSQL, conversation.StatusStreaming, truncated.ID, conversation.StatusTruncated)
		if res.Error != nil {
			return res.Error
		}
		affected = res.RowsAffected

		return tx.Exec(touchOnlySQL, truncated.ConversationID).Error
	})
	if isUniqueViolation(err) {
		return truncated, failure.New(failure.CodeConversationBusy)
	}
	if err != nil {
		return truncated, failure.ErrWithTrace(err)
	}
	if affected != 1 {
		return truncated, failure.New(failure.CodeBadRequest)
	}
	truncated.Status = conversation.StatusStreaming

	return truncated, nil
}

func (r *Repo) Latest(ctx context.Context, conversationID string) (conversation.Message, error) {
	var item conversation.Message
	err := r.db.WithContext(ctx).Where(whereConversationRows, conversationID).Order(orderNewestSeqFirst).Limit(1).Find(&item).Error
	if err != nil {
		return item, failure.ErrWithTrace(err)
	}

	return item, nil
}

func failStale(tx *gorm.DB, conversationID string, staleBefore time.Time) error {
	stale := tx.Model(&conversation.Message{}).Where(whereStaleStreaming, conversationID, conversation.StatusStreaming, staleBefore)

	return stale.Update(columnStatus, conversation.StatusFailed).Error
}

func lastSeq(tx *gorm.DB, conversationID string) (int64, error) {
	var maxSeq int64
	err := tx.Model(&conversation.Message{}).Select(selectMaxSeq).Where(whereConversationRows, conversationID).Scan(&maxSeq).Error

	return maxSeq, err
}

func (r *Repo) History(ctx context.Context, conversationID string, beforeSeq int64) ([]conversation.Message, error) {
	items := make([]conversation.Message, 0, maxHistoryRows)
	err := r.db.WithContext(ctx).
		Select(columnRole, columnContent).
		Where(whereHistory, conversationID, []conversation.MessageStatus{conversation.StatusCompleted, conversation.StatusTruncated}, beforeSeq).
		Order(orderNewestSeqFirst).
		Limit(maxHistoryRows).
		Find(&items).Error
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}
	slices.Reverse(items)

	return items, nil
}

func (r *Repo) Finish(ctx context.Context, in finishInput) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(finishMessageSQL, in.content, in.status, in.usage.InputTokens, in.usage.OutputTokens, in.costMicroUSD,
			in.usage.Estimated, in.firstTokenMS, in.latencyMS, in.messageID).Error
		if err != nil {
			return err
		}

		return tx.Exec(touchOnlySQL, in.conversationID).Error
	})
	if err != nil {
		return failure.ErrWithTrace(err)
	}

	return nil
}

func (r *Repo) SpentSince(ctx context.Context, since time.Time) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&conversation.Message{}).Select(selectSumCost).Where(whereSpent, conversation.RoleAssistant, since).Scan(&total).Error
	if err != nil {
		return 0, failure.ErrWithTrace(err)
	}

	return total, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == uniqueViolation
}
