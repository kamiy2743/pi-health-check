~/workspace/AGENTS.mdを読む

## 概要
- Raspberry Pi 上で公開しているサイト群のURLを定期的にヘルスチェックするための構成。
- Go 製のワンショット実行バイナリを systemd timer で1分ごとに起動し、異常時は Discord に通知する。
- ローカルURLと外部公開URLの両方を監視できる。

## 仕組み
- 監視対象URLは `/home/kamiy2743/workspace/health-check/.conf` に列挙する。
- 実行時は `DISCORD_WEBHOOK_URL`（`.env`）を読み込み、未設定ならエラー終了する。
- 各URLに対してHTTP GETを実行し、2xx 以外を失敗として収集する（タイムアウト3秒）。
- 直前の失敗URL一覧は `/run/healthcheck/state.json` に保存し、今回の結果と差分を計算する。
- 変化があった場合のみ Discord に通知し、疎通不可が継続中の場合は通知を抑制する。

## 通知仕様
- すべて解消 / 一部解消 / 疎通不可発生 / 疎通不可増加の4パターンでメッセージを出し分ける。
- Embed 形式で「疎通不可」「回復」のURL一覧を表示する。
