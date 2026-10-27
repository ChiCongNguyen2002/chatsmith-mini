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
