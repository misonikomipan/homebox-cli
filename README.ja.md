# homebox-cli

[English](README.md)

[![Go Report Card]
(https://goreportcard.com/badge/github.com/misonikomipan/homebox-cli)](https://goreportcard.com/report/github.com/misonikomipan/homebox-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[Homebox](https://github.com/sysadminsmedia/homebox) 在庫管理システムを操作するための強力で使いやすいコマンドラインインターフェース（CLI）です。

**Homebox v0.26.x に対応**（v0.26 では items と locations が単一の entities リソースに統合された API に合わせて更新しています）。

## 主な機能

- **リソース管理**: アイテム、場所、タグ、メンテナンス、通知設定、テンプレート、エンティティタイプに対する CRUD 操作。
- **カスタムフィールド**: エンティティのカスタムフィールドを完全にサポート (`hb items fields`)。
- **API キー**: API キーの作成・一覧・失効 (`hb auth api-keys`) と、キーを直接利用する `hb auth token hb_...`。
- **ラベルメーカー**: アイテム / 場所 / アセットのラベルを PNG で生成 (`hb labelmaker get`)。
- **柔軟な出力形式**: スクリプト用の `json` または人間が読みやすい `table`（表形式）を選択可能。
- **シェル補完**: Bash、Zsh、Fish、PowerShell をサポート。
- **階層構造の表示**: 場所のツリー構造を表示可能（アイテムの有無も選択可能）。
- **データの可搬性**: 在庫アイテムの CSV 形式でのエクスポート・インポート。

## Homebox v0.26 での変更

Homebox v0.26 では `/v1/items` と `/v1/locations` の API が単一の
[`/v1/entities`](https://github.com/sysadminsmedia/homebox/pull/1414) リソースに統合されました。
この CLI は以下のように対応しています。

| 旧エンドポイント (v0.26 以前)      | v0.26 のエンドポイント                      |
| ---------------------------------- | ------------------------------------------- |
| `GET/POST /v1/items`             | `GET/POST /v1/entities`                   |
| `GET/PUT/DELETE /v1/items/{id}`  | `GET/PUT/PATCH/DELETE /v1/entities/{id}`  |
| `GET /v1/items/{id}/path`        | `GET /v1/entities/{id}/path`              |
| `POST /v1/items/{id}/duplicate`  | `POST /v1/entities/{id}/duplicate`        |
| `GET/POST /v1/items/{id}/maintenance` | `GET/POST /v1/entities/{id}/maintenance` |
| `GET/POST /v1/items/export|import` | `GET/POST /v1/entities/export|import`     |
| `GET/POST /v1/items/{id}/attachments` | `POST /v1/entities/{id}/attachments`（`name` フォーム項目が必須に） |
| `GET /v1/locations`              | `GET /v1/entities?isLocation=true`        |
| `GET /v1/locations/tree`         | `GET /v1/entities/tree`                   |
| `POST /v1/locations`             | `POST /v1/entities`（location エンティティタイプを指定） |
| `GET /v1/currency`               | `GET /v1/currencies`                      |
| `PUT /v1/users/change-password`  | `PUT /v1/users/self/change-password`      |
| `GET/POST/PUT/DELETE /v1/labelmakers` | 廃止 — `GET /v1/labelmaker/{entity|item|location|asset}/{id}` を使用 |
| アイテムのカスタムフィールド CRUD  | `PUT /v1/entities/{id}`（fields 配列）で管理 |

v0.26 で新しくなり、この CLI でも利用できる機能:

- **API キー** — 自動化にログインは不要: `hb auth token hb_...` または `HB_TOKEN` でキーを保存し、`hb auth api-keys` で管理。
- **エンティティタイプ** — `hb entity-types list|create|update|delete`。
- **ユーザー設定** — `hb auth settings`。

## インストール

### ソースからビルド

[Go](https://go.dev/doc/install) 1.21 以降がインストールされていることを確認してください。

```bash
git clone https://github.com/misonikomipan/homebox-cli.git
cd homebox-cli
go build -o hb main.go
mv hb /usr/local/bin/ # オプション: パスの通ったディレクトリに移動
```

## クイックスタート

### 1. エンドポイントの設定

Homebox インスタンスの URL を設定します。

```bash
hb config --endpoint https://homebox.example.com
```

### 2. 認証

メールアドレスとパスワードで認証:

```bash
hb login --email your-email@example.com
```

または v0.26 の API キーを使用（自動化におすすめ）:

```bash
hb auth token hb_your_api_key_here
# または
export HB_TOKEN=hb_your_api_key_here
```

### 3. 基本的なコマンド

```bash
# 表形式でアイテム一覧を表示
hb items list --format table

# アイテムを検索
hb items list --query "laptop" --format table

# 場所のツリー構造を表示
hb locations tree --with-items

# アイテムにカスタムフィールドを追加
hb items fields add <item-id> --label "シリアル番号" --value "XYZ-123"

# アイテムのラベル PNG を生成
hb labelmaker get <item-id> -o label.png

# シェル補完スクリプトの生成
hb completion zsh > ~/.zshrc.d/_hb
```

## 使い方

各コマンドの詳細なヘルプは `--help` フラグで確認できます。

```bash
hb --help
hb items --help
hb items create --help
```

## 設定

設定ファイルは `~/.config/hb/config.json` に保存されます。

環境変数を使用して設定を上書きすることも可能です。
- `HB_ENDPOINT`: API エンドポイント URL
- `HB_TOKEN`: 認証トークン（セッショントークンまたは `hb_` API キー）
- `HB_FORMAT`: デフォルトの出力形式 (`json` または `table`)

## 開発

### Git Hooks

コード品質を維持するために、pre-commit と pre-push フックを使用しています。

```bash
# クローン後、以下のコマンドでフックを有効化できます：
chmod +x scripts/hooks/*
git config core.hooksPath scripts/hooks
```

## ライセンス

このプロジェクトは MIT ライセンスの下で提供されています。詳細は [LICENSE](LICENSE) ファイルをご覧ください。
