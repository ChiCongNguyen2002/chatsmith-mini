-- Contract: xoá view tên cũ. Chỉ chạy khi không còn binary nào đọc tên bảng cũ.
BEGIN;

DROP VIEW models;
DROP VIEW sessions;
DROP VIEW user_identities;

COMMIT;
