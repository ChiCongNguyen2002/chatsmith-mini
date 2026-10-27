# Deploy lên GCP

Mục tiêu: một service Cloud Run phục vụ cả API và giao diện, dữ liệu trên Cloud SQL Postgres, khoá bí mật trong Secret Manager. Chi phí được giữ thấp: Cloud Run tối thiểu 0 instance, Cloud SQL chỉ bật khi test.

Các lệnh dưới đây **chưa được chạy trên GCP** (xem trạng thái trong [testing.md](testing.md)). Thay các giá trị trong `<>`.

## 0. Biến dùng chung

```bash
export PROJECT=chatsmith-mini
export REGION=asia-southeast1
export REPO=chatsmith
export DB_INSTANCE=chatsmith-db
export SERVICE=chatsmith-mini
export SA=chatsmith-run@${PROJECT}.iam.gserviceaccount.com
gcloud config set project ${PROJECT}
```

## 1. Ngân sách

Trong Console: **Billing → Budgets & alerts → Create budget**. Đặt ngưỡng 25.000 VND (theo đơn vị tiền tệ của tài khoản) và cảnh báo ở 50%, 90%, 100%. Ngân sách chỉ cảnh báo, không tự tắt dịch vụ.

## 2. Bật API

```bash
gcloud services enable run.googleapis.com sqladmin.googleapis.com \
  secretmanager.googleapis.com artifactregistry.googleapis.com \
  cloudbuild.googleapis.com iam.googleapis.com
```

## 3. Artifact Registry

```bash
gcloud artifacts repositories create ${REPO} --repository-format=docker --location=${REGION}
```

## 4. Cloud SQL

Chọn cấu hình nhỏ nhất để tiết kiệm (kiểm tra giá trên trang giá chính thức trước khi tạo).

```bash
gcloud sql instances create ${DB_INSTANCE} --database-version=POSTGRES_16 \
  --tier=db-f1-micro --region=${REGION}
gcloud sql databases create chatsmith --instance=${DB_INSTANCE}
gcloud sql users create chatsmith --instance=${DB_INSTANCE} --password=<mật khẩu DB>
```

## 5. Secret Manager

Nhập giá trị qua stdin để không lưu vào lịch sử shell:

```bash
printf '%s' '<mật khẩu DB>' | gcloud secrets create db-password --data-file=-
printf '%s' '<khoá Anthropic>' | gcloud secrets create anthropic-api-key --data-file=-
printf '%s' '<client secret Google>' | gcloud secrets create google-client-secret --data-file=-
openssl rand -hex 32 | tr -d '\n' | gcloud secrets create session-secret --data-file=-
```

DATABASE_URL dùng socket của Cloud SQL, không đặt trong repo:
`postgres://chatsmith:<mật khẩu>@/chatsmith?host=/cloudsql/${PROJECT}:${REGION}:${DB_INSTANCE}&sslmode=disable`

## 6. Service account

```bash
gcloud iam service-accounts create chatsmith-run
for role in roles/cloudsql.client roles/secretmanager.secretAccessor; do
  gcloud projects add-iam-policy-binding ${PROJECT} --member=serviceAccount:${SA} --role=${role}
done
```

## 7. OAuth client Google

Trong Console: **APIs & Services → OAuth consent screen**: chọn External, điền tên ứng dụng và email hỗ trợ. Sau đó **Credentials → Create credentials → OAuth client ID → Web application**.

Redirect URI là URL của Cloud Run cộng `/auth/google/callback`. Lần đầu chưa có URL thì thêm sau khi deploy (bước 9).

## 8. Build và đẩy image

Chạy từ thư mục gốc repo (Dockerfile nằm ở thư mục gốc, context là thư mục gốc):

```bash
gcloud builds submit --tag ${REGION}-docker.pkg.dev/${PROJECT}/${REPO}/${SERVICE}:latest \
  --gcs-source-staging-dir=gs://${PROJECT}_cloudbuild/source .
```

Nếu `--gcs-source-staging-dir` báo lỗi bucket, dùng `docker build -f Dockerfile` và `docker push` từ máy.

## 9. Migration

Migration chạy một lần, bằng Cloud SQL Auth Proxy trên máy:

```bash
cloud-sql-proxy ${PROJECT}:${REGION}:${DB_INSTANCE} --port 5433 &
for f in backend/migrations/*.sql; do psql "postgres://chatsmith:<mật khẩu>@127.0.0.1:5433/chatsmith" -f "$f"; done
```

Chỉ chạy lại một file khi biết chắc nó chưa áp dụng. Các file không có `IF NOT EXISTS`.

## 10. Deploy Cloud Run (GUI)

Form **Create service** (Cloud Run → Create service → Continuously deploy from a repository):

- **Region**: `asia-southeast1`. Authentication: **Allow public access**.
- **Variables & Secrets → Environment variables** (chỉ biến thường, không có secret):
  - `GCP_PROJECT_ID` = `chatsmith-mini`
  - `AUTH_GOOGLE_CLIENT_ID` = Client ID của OAuth client
  - `AUTH_PUBLIC_BASE_URL` = `https://example.com` (tạm, sửa sau khi có URL)
  - `AUTH_COOKIE_SECURE` = `true`
- Không bấm **Reference a secret**. App tự đọc `database-url`, `anthropic-api-key`, `google-client-secret`, `session-secret` từ Secret Manager khi khởi động.
- **Cloud SQL connections**: `chatsmith-db`.
- **Security → Service account**: `chatsmith-run` (phải có quyền Secret Manager Secret Accessor và Cloud SQL Client).
- **Resources**: CPU 1, memory 512 MiB. **Request timeout**: 300 giây.
- **Revision scaling**: minimum instances 0, maximum 2.

Sau lần deploy đầu, lấy URL của service và làm bước 11.

## 11. Kiểm tra

```bash
curl <URL>/healthz
```

Mở `<URL>` trên trình duyệt, đăng nhập bằng tài khoản Google, gửi một câu chat. Kiểm tra dòng trong database qua Cloud SQL Auth Proxy.

## Tắt và bật Cloud SQL để tiết kiệm

```bash
gcloud sql instances patch ${DB_INSTANCE} --activation-policy=NEVER    # tắt
gcloud sql instances patch ${DB_INSTANCE} --activation-policy=ALWAYS   # bật
```

Khi Cloud SQL tắt, `/healthz` trả 503 và API trả lỗi. Đây là hành vi đúng, dữ liệu không mất.
