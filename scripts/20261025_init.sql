CREATE TABLE IF NOT EXISTS users (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS conversations (
    id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title            text NOT NULL DEFAULT '',
    last_seq         bigint NOT NULL DEFAULT 0,
    generating_until timestamptz,
    created_at       timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_conversations_user_created
    ON conversations (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS messages (
    id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id   uuid NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
    user_id           uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    seq               bigint NOT NULL,
    role              text NOT NULL CHECK (role IN ('user', 'assistant')),
    content           text NOT NULL,
    status            text NOT NULL CHECK (status IN ('streaming', 'completed', 'failed', 'canceled')),
    client_request_id text,
    created_at        timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT uq_messages_conversation_seq UNIQUE (conversation_id, seq),
    CONSTRAINT uq_messages_conversation_client_request UNIQUE (conversation_id, client_request_id)
);
