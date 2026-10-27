# Chạy nhanh

## Cần có

- Docker và Docker Compose.
- Một OAuth client Google kiểu Web application (xem bước 2 bên dưới).
- Khoá API Anthropic nếu muốn chat thật.
- Go 1.27 nếu muốn chạy không qua Docker.

## 1. Tạo file `.env`

Tạo file `.env` ở thư mục gốc:

```
POSTGRES_PASSWORD=change-me
ANTHROPIC_API_KEY=<khoá Anthropic>
AUTH_GOOGLE_CLIENT_ID=<client id>
AUTH_GOOGLE_CLIENT_SECRET=<client secret>
AUTH_SESSION_SECRET=<chuỗi ngẫu nhiên từ 32 ký tự>
AUTH_PUBLIC_BASE_URL=http://127.0.0.1:8080
AUTH_COOKIE_SECURE=false
```

- `POSTGRES_PASSWORD`: chỉ dùng chữ, số, `-`, `_`, `.`.
- `AUTH_SESSION_SECRET`: tạo bằng `openssl rand -hex 32`.

## 2. Tạo OAuth client trên Google

Trong Google Cloud Console, vào APIs & Services → Credentials → Create credentials → OAuth client ID, chọn **Web application**. Thêm redirect URI:

```
http://127.0.0.1:8080/auth/google/callback
```

Lấy client ID và client secret điền vào `.env`.

## 3. Chạy

```bash
make up
```

Mở http://127.0.0.1:8080 và bấm đăng nhập bằng Google. Kiểm tra máy chủ:

```bash
curl http://127.0.0.1:8080/healthz
```

Kết quả đúng: `{"status":200,"data":"ok"}`.

Dừng: `make down`.

## Chạy không qua Docker

Cần Postgres đang chạy và các biến trong `.env` có trong môi trường:

```bash
export DATABASE_URL='postgres://chatsmith:change-me@127.0.0.1:5432/chatsmith?sslmode=disable'
make run
```

`make run` chạy từ `backend/`, nên `frontend/` nằm ở `../frontend` (đã là mặc định).

## Lưu ý

- Migration chỉ chạy khi volume Postgres mới được tạo. Với volume đã có, phải chạy file trong `backend/migrations/` bằng tay, đúng thứ tự tên file.
- Trang chủ khi chưa đăng nhập hiện nút Google. Trình duyệt có thể báo một lỗi 401 trong console: đó là lời gọi `/v1/me` khi chưa có phiên, không phải lỗi.
