# Chạy nhanh

## Cần có

- Docker và Docker Compose.
- Khoá API Anthropic (không bắt buộc để mở giao diện, nhưng cần để chat thật).
- Go 1.27 nếu muốn chạy không qua Docker.

## Chạy bằng Docker

Tạo file `.env` ở thư mục gốc với hai dòng sau:

```
POSTGRES_PASSWORD=change-me
ANTHROPIC_API_KEY=<khoá của bạn>
```

- `POSTGRES_PASSWORD`: chỉ dùng chữ, số, `-`, `_`, `.` (mật khẩu nằm trong `DATABASE_URL`).
- `ANTHROPIC_API_KEY`: để trống thì không chat được với Claude.

Các khoá khác là tuỳ chọn, xem [configuration.md](configuration.md).

Rồi chạy:

```bash
make up
```

Mở http://127.0.0.1:8080. Kiểm tra máy chủ:

```bash
curl http://127.0.0.1:8080/healthz
```

Kết quả đúng: `{"status":200,"data":"ok"}`.

Dừng:

```bash
make down
```

## Chạy không qua Docker

Cần một Postgres đang chạy và `DATABASE_URL` trong môi trường:

```bash
export DATABASE_URL='postgres://chatsmith:change-me@127.0.0.1:5432/chatsmith?sslmode=disable'
make run
```

`make run` chạy từ thư mục `backend/`, nên `frontend/` nằm ở `../frontend` (đã là mặc định).

## Lưu ý

- Migration chỉ chạy khi volume Postgres mới được tạo. Với volume đã có, phải chạy file trong `backend/migrations/` bằng tay.
- Giao diện đang dùng một ID người dùng tự sinh và lưu trong trình duyệt. Chưa có đăng nhập.
