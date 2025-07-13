# ビルドステージ
FROM golang:1.23-alpine AS builder

# 作業ディレクトリを設定
WORKDIR /app

# 依存関係をコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 実行ステージ
FROM alpine:latest

# ca-certificatesをインストール（HTTPS通信のため）
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# ビルドステージからバイナリをコピー
COPY --from=builder /app/main .

# アプリケーションを実行
CMD ["./main"] 