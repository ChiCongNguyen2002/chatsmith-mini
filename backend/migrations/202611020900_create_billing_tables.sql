-- Billing: gói, hạn mức theo nhóm model, gói đang dùng, lượt dùng theo ngày.
BEGIN;

CREATE TABLE billing_plans (
    id           uuid PRIMARY KEY DEFAULT uuidv7(),
    code         text NOT NULL UNIQUE,            -- 'free', 'pro'
    name         jsonb NOT NULL,                  -- {"en": "Free", "vi": "Miễn phí"}
    daily_limits jsonb NOT NULL CHECK (jsonb_typeof(daily_limits) = 'object'), -- {"basic": 20, "advanced": 5}; giá trị JSON null = không giới hạn (không phải SQL NULL); key thiếu = không được dùng
    is_default   boolean NOT NULL DEFAULT false,  -- Gói cho user chưa mua gì; đúng một gói
    order_index  int NOT NULL DEFAULT 0,
    is_active    boolean NOT NULL DEFAULT true,
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX uq_billing_plans_is_default ON billing_plans (is_default) WHERE is_default;

CREATE TABLE billing_subscriptions (
    id         uuid PRIMARY KEY DEFAULT uuidv7(),
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    plan_id    uuid NOT NULL REFERENCES billing_plans (id),
    status     text NOT NULL CHECK (status IN ('active', 'canceled')), -- 'active': đang dùng, 'canceled': đã huỷ
    expires_at timestamptz NOT NULL,              -- Hết hạn thì về gói mặc định, không cần cron
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX uq_billing_subscriptions_active_user ON billing_subscriptions (user_id) WHERE status = 'active';

CREATE TABLE billing_daily_usages (
    user_id        uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    usage_date     date NOT NULL,                 -- Ngày theo UTC
    quota_group    text NOT NULL,                 -- 'basic', 'advanced': khớp catalog_models.quota_group
    used_count     int NOT NULL DEFAULT 0 CHECK (used_count >= 0), -- Lượt đã dùng trong ngày
    cost_micro_usd bigint NOT NULL DEFAULT 0 CHECK (cost_micro_usd >= 0),
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, usage_date, quota_group)
);
CREATE INDEX idx_billing_daily_usages_usage_date ON billing_daily_usages (usage_date);

ALTER TABLE catalog_models ADD COLUMN quota_group text;
UPDATE catalog_models SET quota_group = 'basic';
UPDATE catalog_models SET quota_group = 'advanced' WHERE code IN ('claude-sonnet-5-5', 'gpt-6.1-sol');
ALTER TABLE catalog_models ALTER COLUMN quota_group SET NOT NULL;

-- FK thêm NOT VALID rồi VALIDATE: không quét toàn bảng messages trong lúc giữ khoá ghi.
ALTER TABLE messages ADD COLUMN plan_id uuid;             -- Gói lúc gửi tin; tin cũ NULL
ALTER TABLE messages ADD CONSTRAINT messages_plan_id_fkey FOREIGN KEY (plan_id) REFERENCES billing_plans (id) NOT VALID;
ALTER TABLE messages VALIDATE CONSTRAINT messages_plan_id_fkey;

INSERT INTO billing_plans (code, name, daily_limits, is_default, order_index) VALUES
    ('free', '{"en": "Free", "vi": "Miễn phí"}', '{"basic": 20, "advanced": 5}', true, 0),
    ('pro',  '{"en": "Pro", "vi": "Pro"}',       '{"basic": 1000, "advanced": 300}', false, 1);

COMMIT;
