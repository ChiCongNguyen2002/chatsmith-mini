CREATE TABLE models (
    id                        text PRIMARY KEY,
    provider                  text NOT NULL CHECK (provider IN ('openai', 'anthropic')),
    display_name              text NOT NULL,
    description               text NOT NULL DEFAULT '',
    input_micro_usd_per_mtok  bigint NOT NULL CHECK (input_micro_usd_per_mtok >= 0),
    output_micro_usd_per_mtok bigint NOT NULL CHECK (output_micro_usd_per_mtok >= 0),
    enabled                   boolean NOT NULL DEFAULT true,
    is_new                    boolean NOT NULL DEFAULT false,
    sort_order                int NOT NULL DEFAULT 0
);

INSERT INTO models (id, provider, display_name, description, input_micro_usd_per_mtok, output_micro_usd_per_mtok, is_new, sort_order) VALUES
    ('gpt-6-luna',        'openai',    'GPT-6 Luna',        'Helps you think creatively and write with depth',          100000,   500000, false, 10),
    ('gpt-6.1-sol',       'openai',    'GPT-6.1 Sol',       'Boosts your productivity by handling routine tasks',      2000000, 10000000, true,  11),
    ('claude-haiku-5-5',  'anthropic', 'Claude Haiku 5.5',  'Fastest model for high-volume, latency-sensitive tasks',   100000,   500000, true,  20),
    ('claude-sonnet-5-5', 'anthropic', 'Claude Sonnet 5.5', 'Speed and intelligence, combined in one model',           2000000, 10000000, false, 21);
