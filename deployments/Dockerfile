FROM golang:1.23.1-alpine AS builder
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN GOOS=linux CGO_ENABLED=0 go build -o insta-archiver ./cmd/web/web.go

FROM scratch
WORKDIR /app

COPY --from=builder /app /app/

CMD [ "/app/insta-archiver" ]
