package chat

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"time"

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

func (r *Repo) BeginTurn(ctx context.Context, conversationID, content, title, model string, staleBefore time.Time) (Message, Message, error) {
	user := Message{ConversationID: conversationID, Role: RoleUser, Content: content, Status: StatusCompleted}
	assistant := Message{ConversationID: conversationID, Role: RoleAssistant, Status: StatusStreaming, Model: &model}

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

		return tx.Exec("UPDATE conversations SET updated_at = now(), title = CASE WHEN title = '' THEN ? ELSE title END WHERE id = ?",
			title, conversationID).Error
	})
	if isUniqueViolation(err) {
		return user, assistant, failure.New(failure.CodeConversationBusy)
	}
	if err != nil {
		return user, assistant, failure.ErrWithTrace(err)
	}

	return user, assistant, nil
}

func (r *Repo) BeginContinue(ctx context.Context, truncated Message, staleBefore time.Time) (Message, error) {
	var affected int64
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := failStale(tx, truncated.ConversationID, staleBefore); err != nil {
			return err
		}
		res := tx.Exec("UPDATE messages SET status = ? WHERE id = ? AND status = ?", StatusStreaming, truncated.ID, StatusTruncated)
		if res.Error != nil {
			return res.Error
		}
		affected = res.RowsAffected

		return tx.Exec("UPDATE conversations SET updated_at = now() WHERE id = ?", truncated.ConversationID).Error
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
	truncated.Status = StatusStreaming

	return truncated, nil
}

func (r *Repo) LatestMessage(ctx context.Context, conversationID string) (Message, error) {
	var item Message
	err := r.db.WithContext(ctx).Where("conversation_id = ?", conversationID).Order("seq DESC").Limit(1).Find(&item).Error
	if err != nil {
		return item, failure.ErrWithTrace(err)
	}

	return item, nil
}

func (r *Repo) History(ctx context.Context, conversationID string, beforeSeq int64) ([]Message, error) {
	items := make([]Message, 0, maxHistoryRows)
	err := r.db.WithContext(ctx).
		Select("role", "content").
		Where("conversation_id = ? AND status IN ? AND seq < ?", conversationID, []MessageStatus{StatusCompleted, StatusTruncated}, beforeSeq).
		Order("seq DESC").
		Limit(maxHistoryRows).
		Find(&items).Error
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}
	slices.Reverse(items)

	return items, nil
}

func (r *Repo) ListMessages(ctx context.Context, conversationID string, afterSeq int64, limit int) ([]Message, error) {
	items := make([]Message, 0, limit)
	err := r.db.WithContext(ctx).Where("conversation_id = ? AND seq > ?", conversationID, afterSeq).Order("seq ASC").Limit(limit).Find(&items).Error
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return items, nil
}

func (r *Repo) FinishMessage(ctx context.Context, in FinishInput) error {
	sources, err := marshalSources(in.Sources)
	if err != nil {
		return failure.ErrWithTrace(err)
	}
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(`UPDATE messages
SET content = ?, status = ?, input_tokens = ?, output_tokens = ?, cost_micro_usd = ?,
    usage_estimated = ?, first_token_ms = ?, latency_ms = ?, sources = COALESCE(?, sources)
WHERE id = ?`, in.Content, in.Status, in.Usage.InputTokens, in.Usage.OutputTokens, in.CostMicroUSD,
			in.Usage.Estimated, in.FirstTokenMS, in.LatencyMS, sources, in.MessageID).Error
		if err != nil {
			return err
		}

		return tx.Exec("UPDATE conversations SET updated_at = now() WHERE id = ?", in.ConversationID).Error
	})
	if err != nil {
		return failure.ErrWithTrace(err)
	}

	return nil
}

func (r *Repo) SpentSince(ctx context.Context, since time.Time) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&Message{}).Select("COALESCE(SUM(cost_micro_usd), 0)").Where("role = ? AND created_at >= ?", RoleAssistant, since).Scan(&total).Error
	if err != nil {
		return 0, failure.ErrWithTrace(err)
	}

	return total, nil
}

func failStale(tx *gorm.DB, conversationID string, staleBefore time.Time) error {
	stale := tx.Model(&Message{}).Where("conversation_id = ? AND status = ? AND created_at < ?", conversationID, StatusStreaming, staleBefore)

	return stale.Update("status", StatusFailed).Error
}

func lastSeq(tx *gorm.DB, conversationID string) (int64, error) {
	var maxSeq int64
	err := tx.Model(&Message{}).Select("COALESCE(MAX(seq), 0)").Where("conversation_id = ?", conversationID).Scan(&maxSeq).Error

	return maxSeq, err
}

func marshalSources(items []Source) (any, error) {
	if len(items) == 0 {
		return nil, nil
	}
	converted := make([]SourceJSON, 0, len(items))
	for _, item := range items {
		converted = append(converted, SourceJSON{URL: item.URL, Title: item.Title})
	}
	raw, err := json.Marshal(converted)
	if err != nil {
		return nil, err
	}

	return string(raw), nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == uniqueViolation
}
