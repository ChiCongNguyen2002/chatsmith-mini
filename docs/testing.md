# Kiểm thử

## Gate Go

Chạy trong `backend/`:

```bash
gofmt -l .        # phải rỗng
go vet ./...
go build ./...
```

Hoặc từ thư mục gốc: `make test`.

Repo không có file `*_test.go`. Kiểm tra logic làm bằng test tạm, chạy xong thì xoá.

## Kiểm tra giao diện

`tests/ui/smoke.js` có hai chế độ:

- **Không có cookie:** kiểm tra màn đăng nhập hiện nút Google.
- **Có cookie `chatsmith_session`:** kiểm tra tên người dùng, gửi một câu, nhận trả lời.

Lấy cookie: đăng nhập trên trình duyệt, mở DevTools → Application → Cookies, sao chép giá trị `chatsmith_session`. Không dán cookie vào tài liệu hay commit.

```bash
docker run --rm --add-host=host.docker.internal:host-gateway -e BASE=http://host.docker.internal:8080 \
  -e SESSION_COOKIE=<giá trị cookie> -v "$PWD/tests:/tests" mcr.microsoft.com/playwright:v1.49.0-jammy \
  sh -c 'mkdir -p /tmp/pw && cd /tmp/pw && npm init -y >/dev/null && npm i playwright@1.49.0 >/dev/null 2>&1 && NODE_PATH=/tmp/pw/node_modules node /tests/ui/smoke.js'
```

Kết quả đúng: có `user:`, `reply:`, `chip:`, và `errors: []`.

## Kiểm tra backend bằng tay

```bash
curl http://127.0.0.1:8080/healthz
curl -i http://127.0.0.1:8080/auth/google/login          # 302 sang accounts.google.com
curl -s -o /dev/null -w '%{http_code}\n' http://127.0.0.1:8080/v1/me   # 401 khi chưa có phiên
```

## Đã kiểm chứng và còn lại

- Đã kiểm (CM-22, trên stack local, với phiên tạo bằng cùng secret): `/v1/me` đúng tài khoản; cookie bị sửa bị từ chối; chat thật qua cookie hoàn thành; hai tài khoản không thấy chat của nhau; trang chủ không có phiên hiện nút đăng nhập.
- Chưa kiểm: đăng nhập Google thật (chưa có OAuth client trên GCP), ngân sách và trần chi tiêu sau thay đổi, `/healthz` khi DB tắt, nút Stop, OpenAI, Cloud Run.
