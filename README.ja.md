# hb — Homebox CLI

[Homebox](https://homebox.software) REST API の Go 製コマンドラインツールです。

## インストール

### ソースからビルド

```bash
git clone https://github.com/misonikomipan/homebox-cli.git
cd homebox-cli
go build -o hb .
mv hb /usr/local/bin/hb
```

## クイックスタート

```bash
# 1. エンドポイントを設定（デフォルト: https://homebox.mizobuchi.dev）
hb config --endpoint https://your-homebox.example.com

# 2. ログイン
hb login --email you@example.com

# 3. 接続確認
hb status

# 4. アイテム一覧
hb items list
```

## 設定

設定は `~/.config/hb/config.json`（パーミッション `0600`）に保存されます。

| キー       | 説明                   |
|------------|------------------------|
| `endpoint` | Homebox サーバーの URL |
| `token`    | Bearer 認証トークン    |

### 環境変数

環境変数は設定ファイルより優先されます。

| 変数名        | 説明                         |
|---------------|------------------------------|
| `HB_ENDPOINT` | API エンドポイントの上書き   |
| `HB_TOKEN`    | 認証トークンの上書き         |

```bash
HB_ENDPOINT=https://homebox.example.com hb items list
```

## コマンド一覧

```
hb [command]

トップレベル:
  login           ログインしてトークンを保存
  logout          ログアウトしてトークンを削除
  status          API ステータスを取得
  config          設定の確認・変更
  guide           使用例のクイックリファレンスを表示
  currency        通貨情報を取得
  barcode-search  バーコード/EAN で商品を検索
```

### auth（認証・アカウント）

```bash
hb auth me                            # 現在のユーザー情報
hb auth refresh                       # トークンを更新
hb auth update-me --name "新しい名前" # プロフィール更新
hb auth change-password               # パスワード変更
```

### items（アイテム管理）

```bash
hb items list                                       # 全アイテム一覧
hb items list --query "laptop" --page-size 20       # 検索
hb items list --location <id>                       # ロケーションで絞り込み
hb items list --label <id>                          # タグで絞り込み
hb items get <id>
hb items create --name "MacBook Pro" --location <id>
hb items create --name "カメラ" --quantity 1 --purchase-price 80000 --notes "Sony A7"
hb items update <id> --name "新しい名前"
hb items delete <id> --yes
hb items duplicate <id>                             # 複製
hb items path <id>                                  # 階層パスを表示
hb items maintenance <id>                           # メンテナンスログ
hb items export --output items.csv
hb items import items.csv
hb items asset <asset-id>                           # アセットIDで検索
hb items attachments upload <item-id> photo.jpg     # 添付ファイルをアップロード
hb items attachments delete <item-id> <attachment-id>
```

### locations（ロケーション管理）

```bash
hb locations list
hb locations tree
hb locations tree --with-items          # アイテムを含むツリー
hb locations get <id>
hb locations create --name "書斎"
hb locations create --name "棚A" --parent <親のid>
hb locations update <id> --name "新名前"
hb locations delete <id> --yes
```

### tags（タグ管理）

```bash
hb tags list
hb tags get <id>
hb tags create --name "電子機器" --color "#3b82f6"
hb tags update <id> --name "ガジェット"
hb tags delete <id> --yes
```

### groups（グループ管理）

```bash
hb groups info
hb groups stats                                  # 統計情報
hb groups members                                # メンバー一覧
hb groups update --name "自宅" --currency "JPY"
hb groups invite --uses 3 --expiry-days 7        # 招待リンク作成
```

### maintenance（メンテナンス管理）

```bash
hb maintenance list
hb maintenance create --item <id> --name "オイル交換" --cost 3000
hb maintenance update <id> --completed-date 2026-03-12
hb maintenance delete <id> --yes
```

### notifiers（通知管理）

```bash
hb notifiers list
hb notifiers create --name "Slack" --url https://hooks.slack.com/...
hb notifiers update <id> --active=false
hb notifiers test                                # テスト通知を送信
hb notifiers delete <id> --yes
```

### templates（テンプレート管理）

```bash
hb templates list
hb templates get <id>
hb templates create --name "PCテンプレート"
hb templates update <id> --name "ノートPCテンプレート"
hb templates create-item <template-id> --location <id>  # テンプレートからアイテム作成
hb templates delete <id> --yes
```

## Tips

```bash
# jq で整形・フィルタリング
hb items list | jq '.items[]?.name'

# ロケーションのIDを一覧表示
hb locations list | jq '.[].id'

# サブコマンドのヘルプ
hb items --help
hb items create --help

# クイックリファレンスを表示
hb guide
```
