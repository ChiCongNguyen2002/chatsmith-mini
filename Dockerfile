FROM golang:1.25 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/app ./cmd/app

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/app ./app
COPY scripts ./scripts
ENV GIN_MODE=release
EXPOSE 8080
ENTRYPOINT ["/app/app"]
