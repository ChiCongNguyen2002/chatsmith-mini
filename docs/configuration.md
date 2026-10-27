# Cấu hình

Mỗi khoá trong `backend/configs/default.go` có giá trị mặc định. Có thể ghi đè bằng biến môi trường, theo quy tắc: `SERVER.PORT` thành `SERVER_PORT`, `CHAT.DAILY_BUDGET_USD` thành `CHAT_DAILY_BUDGET_USD`.

File cấu hình mặc định: `backend/configs/dev.config.yaml`. Đổi file bằng `CONFIG_FILE`.

## Các khoá quan trọng

| Khoá | Mặc định | Ý nghĩa |
|---|---|---|
| `SERVER.PORT` | `8080` | Cổng HTTP |
| `SERVER.FRONTEND_DIR` | `../frontend` | Thư mục giao diện. Docker đặt `/app/frontend` |
| `DATABASE.URL` | (bắt buộc) | Lấy từ `DATABASE_URL` |
| `PROVIDER.ANTHROPIC_API_KEY` | (trống) | Khoá Anthropic. Trống thì không dùng được model Claude |
| `PROVIDER.ANTHROPIC_EFFORT` | `low` (mặc định code), `medium` (file dev) | Mức suy nghĩ của Claude |
| `CHAT.DEFAULT_MODEL` | `claude-haiku-5-5` | Model mặc định |
| `CHAT.MAX_OUTPUT_TOKENS` | `2048` | Trần token ra cho mỗi lượt |
| `CHAT.DAILY_BUDGET_USD` | `2.0` | Trần chi tiêu mỗi ngày cho cả hệ thống. Vượt thì trả lỗi 429 |
| `CHAT.TIMEOUT` | `120s` | Thời gian tối đa cho một lượt trả lời |

## Lưu ý

- `PROVIDER.ANTHROPIC_EFFORT` có hai giá trị khác nhau: mặc định trong code là `low`, file dev là `medium`. Với `low`, Haiku có thể bỏ suy nghĩ và trả lời sai bài toán suy luận đã đo. Nên đặt `medium` khi chạy thật.

## Đăng nhập Google (AUTH)

| Biến môi trường | Mặc định | Ý nghĩa |
|---|---|---|
| `GCP_PROJECT_ID` | (trống) | Khi có giá trị, app đọc `database-url`, `anthropic-api-key`, `google-client-secret`, `session-secret` từ Secret Manager của project này (phiên bản `latest`) và bỏ qua biến môi trường tương ứng. Dùng trên Cloud Run; local để trống |
| `AUTH_GOOGLE_CLIENT_ID` | (bắt buộc) | Client ID của OAuth client kiểu Web application |
| `AUTH_GOOGLE_CLIENT_SECRET` | (bắt buộc) | Client secret của OAuth client. Không đưa vào repo |
| `AUTH_PUBLIC_BASE_URL` | (bắt buộc) | URL công khai của app. Redirect URI là `<URL>/auth/google/callback` |
| `AUTH_SESSION_SECRET` | (bắt buộc, từ 32 ký tự) | Khoá ký cookie phiên. Tạo bằng `openssl rand -hex 32` |
| `AUTH_SESSION_TTL` | `168h` | Thời hạn cookie phiên (7 ngày) |
| `AUTH_STATE_TTL` | `5m` | Thời hạn cookie `state` trong lúc đăng nhập |
| `AUTH_COOKIE_SECURE` | `true` | Đặt `false` khi chạy trên `http://127.0.0.1` |

Thiếu một khoá bắt buộc thì app không khởi động.
