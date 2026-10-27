FROM golang:1.27 AS build
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 go build -o /out/app ./cmd/app

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/app ./app
COPY backend/configs ./configs
COPY frontend ./frontend
ENV GIN_MODE=release SERVER_FRONTEND_DIR=/app/frontend
EXPOSE 8080
ENTRYPOINT ["/app/app"]
