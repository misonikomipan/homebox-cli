# homebox-cli

[![Go Report Card](https://goreportcard.com/badge/github.com/misonikomipan/homebox-cli)](https://goreportcard.com/report/github.com/misonikomipan/homebox-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[Homebox](https://github.com/sysadminsmedia/homebox) 在庫管理システムを操作するための強力で使いやすいコマンドラインインターフェース（CLI）です。

## 主な機能

- **リソース管理**: アイテム、場所、タグ、メンテナンス、通知設定、テンプレートに対する CRUD 操作。
- **カスタムフィールド**: アイテムのカスタムフィールドを完全にサポート (`hb items fields`)。
- **ラベルメーカー**: ラベルメーカー設定の管理。
- **柔軟な出力形式**: スクリプト用の `json` または人間が読みやすい `table`（表形式）を選択可能。
- **シェル補完**: Bash、Zsh、Fish、PowerShell をサポート。
- **階層構造の表示**: 場所のツリー構造を表示可能（アイテムの有無も選択可能）。
- **データの可搬性**: 在庫アイテムの CSV 形式でのエクスポート・インポート。

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

### 2. ログイン

メールアドレスとパスワードで認証します。

```bash
hb login --email your-email@example.com
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
- `HB_TOKEN`: 認証トークン
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
