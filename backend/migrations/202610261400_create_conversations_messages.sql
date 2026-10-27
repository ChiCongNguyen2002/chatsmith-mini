CREATE TABLE conversations (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users (id),
    title      text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_conversations_user_updated ON conversations (user_id, updated_at DESC, id DESC);

CREATE TABLE messages (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id uuid NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
    seq             bigint NOT NULL,
    role            text NOT NULL CHECK (role IN ('user', 'assistant')),
    content         text NOT NULL,
    status          text NOT NULL CHECK (status IN ('streaming', 'completed', 'failed', 'canceled')),
    created_at      timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT uq_messages_conversation_seq UNIQUE (conversation_id, seq)
);
