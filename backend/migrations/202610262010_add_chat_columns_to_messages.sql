ALTER TABLE messages
    ADD COLUMN model            text REFERENCES models (id),
    ADD COLUMN input_tokens     bigint CHECK (input_tokens >= 0),
    ADD COLUMN output_tokens    bigint CHECK (output_tokens >= 0),
    ADD COLUMN cost_micro_usd   bigint CHECK (cost_micro_usd >= 0),
    ADD COLUMN usage_estimated  boolean,
    ADD COLUMN first_token_ms   int,
    ADD COLUMN latency_ms       int;

CREATE UNIQUE INDEX uq_messages_one_streaming ON messages (conversation_id) WHERE status = 'streaming';

CREATE INDEX idx_messages_assistant_created ON messages (created_at) WHERE role = 'assistant';
