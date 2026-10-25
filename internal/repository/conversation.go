package repository

import (
	"context"
	"time"

	"chatsmith-mini/pkg/failure"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConversationItem struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageItem struct {
	ID        string    `json:"id"`
	Seq       int64     `json:"seq"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Conversation struct {
	db *pgxpool.Pool
}

func NewConversation(db *pgxpool.Pool) *Conversation {
	return &Conversation{db: db}
}

const (
	createSQL = `
WITH ensure_user AS (
    INSERT INTO users (id) VALUES ($1) ON CONFLICT (id) DO NOTHING
)
INSERT INTO conversations (user_id, title) VALUES ($1, $2)
RETURNING id, title, created_at`

	listFirstSQL = `
SELECT id, title, created_at FROM conversations
WHERE user_id = $1
ORDER BY created_at DESC, id DESC
LIMIT $2`

	listAfterSQL = `
SELECT id, title, created_at FROM conversations
WHERE user_id = $1 AND created_at <= $3 AND (created_at < $3 OR id < $4)
ORDER BY created_at DESC, id DESC
LIMIT $2`

	ownedSQL = `SELECT EXISTS (SELECT 1 FROM conversations WHERE id = $1 AND user_id = $2)`

	messagesSQL = `
SELECT id, seq, role, content, status, created_at FROM messages
WHERE conversation_id = $1 AND seq > $2
ORDER BY seq ASC
LIMIT $3`
)

func (r *Conversation) Create(ctx context.Context, userID, title string) (*ConversationItem, error) {
	var item ConversationItem
	err := r.db.QueryRow(ctx, createSQL, userID, title).Scan(&item.ID, &item.Title, &item.CreatedAt)
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return &item, nil
}

func (r *Conversation) List(ctx context.Context, userID string, limit int, afterCreatedAt *time.Time, afterID string) ([]ConversationItem, error) {
	var (
		rows pgx.Rows
		err  error
	)
	if afterCreatedAt == nil {
		rows, err = r.db.Query(ctx, listFirstSQL, userID, limit)
	} else {
		rows, err = r.db.Query(ctx, listAfterSQL, userID, limit, *afterCreatedAt, afterID)
	}
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}
	defer rows.Close()

	items := make([]ConversationItem, 0, limit)
	for rows.Next() {
		var item ConversationItem
		if err := rows.Scan(&item.ID, &item.Title, &item.CreatedAt); err != nil {
			return nil, failure.ErrWithTrace(err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return items, nil
}

func (r *Conversation) Owned(ctx context.Context, conversationID, userID string) (bool, error) {
	var owned bool
	if err := r.db.QueryRow(ctx, ownedSQL, conversationID, userID).Scan(&owned); err != nil {
		return false, failure.ErrWithTrace(err)
	}

	return owned, nil
}

func (r *Conversation) Messages(ctx context.Context, conversationID string, afterSeq int64, limit int) ([]MessageItem, error) {
	rows, err := r.db.Query(ctx, messagesSQL, conversationID, afterSeq, limit)
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}
	defer rows.Close()

	items := make([]MessageItem, 0, limit)
	for rows.Next() {
		var item MessageItem
		if err := rows.Scan(&item.ID, &item.Seq, &item.Role, &item.Content, &item.Status, &item.CreatedAt); err != nil {
			return nil, failure.ErrWithTrace(err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return items, nil
}
