# Kiến trúc

## Thành phần

```
trình duyệt ── HTTP / SSE ──▶ backend (Go, gin) ──▶ Postgres
                                   │
                                   └──▶ Anthropic / OpenAI
```

- `frontend/`: HTML, CSS, JS tĩnh. Backend phục vụ tại `/` và `/static/`.
- `backend/`: chia theo lớp. `controller` nhận request, `service` xử lý nghiệp vụ, `repo` truy vấn DB. Mỗi domain nằm trong `backend/internal/<domain>/`.
- `backend/internal/provider/`: adapter cho từng nhà cung cấp, cùng một luồng stream.

## Luồng một lượt chat

1. Trình duyệt gửi `POST /v1/conversations/{id}/messages` kèm model.
2. Server kiểm tra model, quyền sở hữu cuộc chat, và trần chi tiêu trong ngày.
3. Ghi hai tin: tin của user và tin trợ lý ở trạng thái `streaming`.
4. Gọi model và đẩy từng chữ về trình duyệt qua SSE (`message.delta`).
5. Ghi kết quả: nội dung, trạng thái, token, chi phí, thời gian.
6. Gửi `message.completed` hoặc `message.failed`.

Trạng thái cuối: `completed`, `truncated` (hết giới hạn token), `canceled` (người dùng dừng), `failed` (lỗi nhà cung cấp hoặc hết thời gian).

## Dữ liệu

| Bảng | Nội dung |
|---|---|
| `users` | Người dùng (hiện là ID tự sinh) |
| `conversations` | Cuộc chat, có `updated_at` để xếp danh sách |
| `messages` | Tin nhắn, token, chi phí, trạng thái. Chỉ một tin `streaming` cho mỗi cuộc chat |
| `models` | Danh mục model, giá theo micro-USD cho một triệu token |

Chi phí: `ceil((token_vào × giá_vào + token_ra × giá_ra) / 1 000 000)`, đơn vị micro-USD.

## Chưa có

- Đăng nhập Google và phiên đăng nhập thật.
- Gói Free và Pro, giới hạn lượt theo ngày và theo phút.
- Nhiều nhà cung cấp: OpenAI đã có code nhưng chưa gọi được vì hết credit. Gemini, DeepSeek, Grok chưa có.
- Tiếp tục stream khi mất kết nối. Hiện tại mất kết nối là dừng lượt đó.
- Deploy lên GCP.
