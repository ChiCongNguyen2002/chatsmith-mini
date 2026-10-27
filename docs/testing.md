# Kiểm thử

## Gate Go

Chạy trong thư mục `backend/`:

```bash
gofmt -l .        # phải rỗng
go vet ./...
go build ./...
```

Hoặc từ thư mục gốc: `make test`.

Chưa có file `*_test.go` nào trong repo. Các kiểm tra logic hiện làm bằng test tạm, chạy xong thì xoá.

## Kiểm tra giao diện

`tests/ui/smoke.js` mở trang chủ, mở menu chọn model, và gửi một câu. Chạy bằng Playwright trong Docker (không cần cài trình duyệt trên máy):

```bash
docker run --rm --add-host=host.docker.internal:host-gateway -v "$PWD/tests:/tests" \
  mcr.microsoft.com/playwright:v1.49.0-jammy sh -c \
  'mkdir -p /tmp/pw && cd /tmp/pw && npm init -y >/dev/null && npm i playwright@1.49.0 >/dev/null 2>&1 && NODE_PATH=/tmp/pw/node_modules BASE=http://host.docker.internal:8080 node /tests/ui/smoke.js'
```

Kết quả đúng: in ra `errors: []`.

## Kiểm tra backend bằng tay

```bash
curl http://127.0.0.1:8080/healthz
curl -sN -X POST http://127.0.0.1:8080/v1/conversations/<id>/messages \
  -H "X-User-ID: <uuid>" -H 'Content-Type: application/json' \
  -d '{"content":"Say hello in one short sentence.","model":"claude-haiku-5-5"}'
```

Kết quả đúng: có các sự kiện `message.created`, `message.delta`, và kết thúc bằng `message.completed` kèm chi phí.

## Đã kiểm chứng và còn lại

- Đã kiểm: lượt bình thường qua SSE, lượt bị cắt (`truncated`), chi phí khớp công thức, giao diện gửi tin và hiển thị câu trả lời.
- Chưa kiểm: nút Stop trên bản hiện tại, bài toán suy luận với `medium`, OpenAI (hết credit), DB chết giữa lúc stream.
