# Pixela MCP Server

Pixela APIを操作するためのMCP（Model Context Protocol）サーバーです。Go言語で実装されています。

## 機能

このMCPサーバーは以下のPixela API操作をサポートしています：

- **ユーザー作成**: Pixelaで新しいユーザーを作成
- **グラフ作成**: ユーザーのグラフを作成
- **ピクセル投稿**: グラフにデータを投稿

## セットアップ

### 前提条件

- Go 1.21以上
- Pixelaアカウント（https://pixe.la/）

### インストール

#### 方法1: 直接実行

1. リポジトリをクローン
```bash
git clone https://github.com/a-know/pixela-mcp.git
cd pixela-mcp
```

2. 依存関係をインストール
```bash
go mod tidy
```

3. サーバーを起動
```bash
go run .
```

#### 方法2: Dockerを使用

1. リポジトリをクローン
```bash
git clone https://github.com/a-know/pixela-mcp.git
cd pixela-mcp
```

2. Docker Composeで起動
```bash
docker-compose up -d
```

または、Dockerfileから直接ビルド
```bash
docker build -t pixela-mcp .
docker run -p 8080:8080 pixela-mcp
```

デフォルトではポート8080でサーバーが起動します。環境変数`PORT`で変更可能です。

## 使用方法

### MCPクライアントでの設定

#### 直接実行の場合
MCPクライアント（例：Cursor）で以下の設定を追加してください：

```json
{
  "mcpServers": {
    "pixela": {
      "command": "go",
      "args": ["run", "."],
      "cwd": "/path/to/pixela-mcp"
    }
  }
}
```

#### Dockerの場合
```json
{
  "mcpServers": {
    "pixela": {
      "command": "docker",
      "args": ["run", "--rm", "-p", "8080:8080", "pixela-mcp"],
      "cwd": "/path/to/pixela-mcp"
    }
  }
}
```

### 利用可能なツール

#### create_user
Pixelaでユーザーを作成します。

**パラメータ:**
- `username`: ユーザー名
- `token`: 認証トークン
- `agreeTermsOfService`: 利用規約への同意（"yes"/"no"）
- `notMinor`: 未成年でないことの確認（"yes"/"no"）

#### create_graph
ユーザーのグラフを作成します。

**パラメータ:**
- `username`: ユーザー名
- `token`: 認証トークン
- `graphID`: グラフID
- `name`: グラフ名
- `unit`: 単位
- `type`: グラフタイプ（"int"/"float"）
- `color`: グラフの色

#### post_pixel
グラフにピクセルを投稿します。

**パラメータ:**
- `username`: ユーザー名
- `token`: 認証トークン
- `graphID`: グラフID
- `date`: 日付（yyyyMMdd形式、省略時は今日）
- `quantity`: 数量

## 開発

### プロジェクト構造

```
pixela-mcp/
├── main.go              # MCPサーバーのメインエントリーポイント
├── tools.go             # MCPツールの実装
├── main_test.go         # テストファイル
├── pixela/
│   └── client.go        # Pixela APIクライアント
├── go.mod               # Goモジュール定義
├── Dockerfile           # Dockerイメージ定義
├── docker-compose.yml   # Docker Compose設定
├── .dockerignore        # Docker除外ファイル
└── README.md            # このファイル
```

### テスト

```bash
go test -v
```

### Dockerビルド

```bash
# イメージをビルド
docker build -t pixela-mcp .

# コンテナを実行
docker run -p 8080:8080 pixela-mcp

# Docker Composeで起動
docker-compose up -d
```

## ライセンス

MIT License

## 作者

a-know (https://github.com/a-know) 