-- Expand: đổi tên bảng và thêm cột mới. Code cũ vẫn chạy nhờ các view mang tên cũ.
-- View chỉ đọc: code cũ không ghi vào các bảng này. Xoá view ở file *_contract.sql, sau khi code mới đã chạy ổn định.
-- provider_model = code là giá trị backfill. Code mới đọc provider_model để gửi cho provider.
BEGIN;

ALTER TABLE models RENAME TO catalog_models;
ALTER TABLE catalog_models RENAME COLUMN id TO code;
ALTER TABLE catalog_models RENAME COLUMN enabled TO is_active;
ALTER TABLE catalog_models RENAME COLUMN sort_order TO order_index;
ALTER TABLE catalog_models ADD COLUMN provider_model text;
UPDATE catalog_models SET provider_model = code;
ALTER TABLE catalog_models ALTER COLUMN provider_model SET NOT NULL;
ALTER TABLE catalog_models
    ADD COLUMN id         uuid NOT NULL DEFAULT uuidv7(),
    ADD COLUMN created_at timestamptz NOT NULL DEFAULT now(),
    ADD COLUMN updated_at timestamptz NOT NULL DEFAULT now();

-- messages.model trỏ vào code, nên FK phải được tạo lại trên cột mới.
ALTER TABLE messages DROP CONSTRAINT messages_model_fkey;
ALTER TABLE catalog_models ADD CONSTRAINT uq_catalog_models_code UNIQUE (code);
ALTER TABLE catalog_models DROP CONSTRAINT models_pkey;
ALTER TABLE catalog_models ADD CONSTRAINT catalog_models_pkey PRIMARY KEY (id);
ALTER TABLE messages ADD CONSTRAINT messages_model_fkey FOREIGN KEY (model) REFERENCES catalog_models (code);

CREATE VIEW models AS
SELECT code AS id, provider, display_name, description,
       input_micro_usd_per_mtok, output_micro_usd_per_mtok,
       is_active AS enabled, is_new, order_index AS sort_order
FROM catalog_models;

ALTER TABLE sessions RENAME TO auth_sessions;
ALTER INDEX sessions_user_id RENAME TO idx_auth_sessions_user_id;
CREATE VIEW sessions AS SELECT * FROM auth_sessions;

ALTER TABLE user_identities RENAME TO auth_identities;
ALTER INDEX user_identities_verified_email RENAME TO idx_auth_identities_verified_email;
CREATE VIEW user_identities AS SELECT * FROM auth_identities;

COMMIT;
