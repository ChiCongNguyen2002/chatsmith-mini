# Chatsmith Mini

Chat với nhiều model AI (hiện có Claude Haiku 5.5 và Sonnet 5.5). Câu trả lời hiện dần từng chữ (SSE), lịch sử được lưu, và mỗi lượt có chi phí theo bảng giá model.

## Cấu trúc

| Thư mục | Nội dung |
|---|---|
| `backend/` | API Go (gin, GORM, Postgres), migrations, cấu hình |
| `frontend/` | Giao diện HTML, CSS, JS tĩnh, do backend phục vụ |
| `deploy/` | `docker-compose.yml` cho máy dev |
| `docs/` | Hướng dẫn chạy, cấu hình, kiến trúc, kiểm thử |
| `tests/` | Script kiểm tra giao diện bằng Playwright |

## Bắt đầu nhanh

Xem [docs/quickstart.md](docs/quickstart.md).

```bash
make up                # cần file .env, xem docs/quickstart.md; mở http://127.0.0.1:8080
```

## Tài liệu

- [Chạy nhanh](docs/quickstart.md)
- [Cấu hình](docs/configuration.md)
- [Kiến trúc](docs/architecture.md)
- [Kiểm thử](docs/testing.md)
- [Deploy lên GCP](docs/deployment.md)

## Trạng thái

Chat với Claude, lưu lịch sử theo tài khoản Google, tính chi phí và trần chi tiêu mỗi ngày đã chạy trên máy dev. Đăng nhập Google đã có trong code, chưa deploy lên GCP. Chưa có: gói Free/Pro, giới hạn lượt dùng, nhiều nhà cung cấp. Chi tiết trong [docs/architecture.md](docs/architecture.md#chưa-có).
