# Pixela MCP Server

Pixela APIを操作するためのMCP（Model Context Protocol）サーバーです。Go言語で実装されています。

## 機能

このMCPサーバーは以下のPixela API操作をサポートしています：

### ユーザー管理
- **ユーザー作成** (`create_user`): Pixelaで新しいユーザーを作成
- **ユーザー情報更新** (`update_user`): ユーザーの認証トークンやサンクスコードを更新
- **ユーザープロフィール更新** (`update_user_profile`): ユーザーのプロフィール情報を更新
- **ユーザー削除** (`delete_user`): ユーザーを削除

### グラフ管理
- **グラフ作成** (`create_graph`): ユーザーのグラフを作成
- **グラフ定義更新** (`update_graph`): グラフの定義を更新
- **グラフ削除** (`delete_graph`): 特定のグラフを削除
- **グラフ定義一覧取得** (`get_graphs`): ユーザーの全グラフ定義を取得
- **特定グラフ定義取得** (`get_graph_definition`): 特定のグラフ定義を取得

### ピクセル管理
- **ピクセル投稿** (`post_pixel`): グラフにデータを投稿

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

#### ユーザー管理

##### create_user
Pixelaでユーザーを作成します。

**パラメータ:**
- `username`: ユーザー名
- `token`: 認証トークン
- `agreeTermsOfService`: 利用規約への同意（"yes"/"no"）
- `notMinor`: 未成年でないことの確認（"yes"/"no"）

##### update_user
ユーザー情報を更新します。

**パラメータ:**
- `username`: ユーザー名
- `token`: 現在の認証トークン
- `newToken`: 新しい認証トークン
- `thanksCode`: サンクスコード（オプション）

##### update_user_profile
ユーザープロフィールを更新します。

**パラメータ:**
- `username`: ユーザー名
- `token`: 認証トークン
- `displayName`: 表示名（オプション）
- `profileURL`: プロフィールURL（オプション）
- `description`: プロフィール説明（オプション）
- `avatarURL`: アバター画像URL（オプション）
- `twitter`: Twitterユーザー名（オプション）
- `github`: GitHubユーザー名（オプション）
- `website`: ウェブサイトURL（オプション）

##### delete_user
ユーザーを削除します。

**パラメータ:**
- `username`: ユーザー名
- `token`: 認証トークン

#### グラフ管理

##### create_graph
ユーザーのグラフを作成します。

**パラメータ:**
- `username`: ユーザー名
- `token`: 認証トークン
- `graphID`: グラフID
- `name`: グラフ名
- `unit`: 単位
- `type`: グラフタイプ（"int"/"float"）
- `color`: グラフの色

##### update_graph
グラフ定義を更新します。

**パラメータ:**
- `username`: ユーザー名
- `token`: 認証トークン
- `graphID`: グラフID
- `name`: グラフ名（オプション）
- `unit`: 単位（オプション）
- `color`: グラフの色（オプション）
- `timezone`: タイムゾーン（オプション）
- `selfSufficient`: 自己充足（"yes"/"no"、オプション）
- `isSecret`: 秘密グラフ（"yes"/"no"、オプション）
- `publishOptionalData`: オプションデータ公開（"yes"/"no"、オプション）

##### delete_graph
特定のグラフを削除します。

**パラメータ:**
- `username`: ユーザー名
- `token`: 認証トークン
- `graphID`: グラフID

##### get_graphs
ユーザーのグラフ定義一覧を取得します。

**パラメータ:**
- `username`: ユーザー名
- `token`: 認証トークン

##### get_graph_definition
特定のグラフ定義を取得します。

**パラメータ:**
- `username`: ユーザー名
- `token`: 認証トークン
- `graphID`: グラフID

#### ピクセル管理

##### post_pixel
グラフにピクセルを投稿します。

**パラメータ:**
- `username`: ユーザー名
- `token`: 認証トークン
- `graphID`: グラフID
- `date`: 日付（yyyyMMdd形式、省略時は今日）
- `quantity`: 数量

## 技術仕様

### Pixela API対応
- Pixela APIの型揺れ（bool型とstring型の混在）に対応するため、カスタム型`BoolString`を実装
- 各APIエンドポイントに対応するメソッドを実装
- エラーハンドリングとレスポンス解析を適切に実装

### MCPプロトコル対応
- MCP 2024-11-05プロトコルバージョンに対応
- ツールリストの動的更新に対応
- JSON-RPC 2.0形式でのリクエスト・レスポンス処理

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
├── cursor_log.md        # 開発ログ
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

## 注意事項

### Pixela APIの制限
- 一部の機能（プロフィールのSNS連携、グラフ定義取得など）はPixelaサポーター限定
- 通常ユーザーでは25%の確率でリクエストが拒否される場合がある
- グラフ削除は永続的な操作で、取り消しはできない

### 実装上の注意点
- グラフ更新では`type`フィールドは更新不可（作成時のみ設定可能）
- オプションパラメータは指定されたもののみ更新される
- APIレスポンスの型揺れに対応するため、カスタム型を使用

## ライセンス

MIT License

## 作者

a-know (https://github.com/a-know) 