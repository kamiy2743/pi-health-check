~/workspace/AGENTS.mdを読む

## 概要
- このリポジトリは、Raspberry Pi 上で動かしている公開サイトのヘルスチェックを学習・運用するためのもの。
- 主な目的は、`/health` エンドポイントの監視、異常検知、通知（例: Discord Webhook）、定期実行（`systemd timer`）を段階的に実装すること。

## 最終目的
- 複数の監視対象URL（ローカル / 外部公開）を定期的にヘルスチェックできるようにする。
- すべてのURLを最後まで確認したうえで、失敗したURL一覧とエラー内容を取得できるようにする。
- 異常時に Discord Webhook へ通知し、必要に応じて復旧通知も送れるようにする。
- `systemd service` + `systemd timer` により、Raspberry Pi 起動後も継続して自動監視できるようにする。

## 現在までの進捗
- Nginx 側に `/health` エンドポイントを追加し、ローカルおよび外部URLで疎通確認できる状態にした。
- `check.conf` を作成し、監視対象URLを設定ファイルで管理できるようにした。
- `healthcheck.service`（`Type=oneshot`）を作成し、`systemd` 経由で手動実行できることを確認した。
- `healthcheck.timer` を作成し、1分ごとの定期実行と `journalctl` でのログ確認ができることを確認した。
- 障害テストとして Nginx 停止時の失敗（接続拒否 / 外部 `502`）と、復旧後の正常化を確認した。
- Go 実装 `check.go` を追加し、複数URLのヘルスチェックを最後まで実行できるようにした。
- `DISCORD_WEBHOOK_URL` を環境変数から読み込み、未設定時はエラー終了するようにした。
- Discord 通知は Embed 形式にし、URL とエラー内容をフィールドで縦に並べる形式にした。
- `check.sh` と `discord_payload.py` は試行途中の成果物として残っている（Go 版へ移行予定）。
- `healthcheck.env` を使って Webhook URL を分離し、手動実行時は `set -a; source` で読み込む方針。
