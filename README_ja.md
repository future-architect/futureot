# futureot

[OpenCensus](https://opencensus.io/)ヘルパー集です。

将来、[OpenTelemetry](https://opentelemetry.io/)もサポートする予定です。

```sh
$ go get github.com/future-architect/futureot/...
```

## 設定ヘルパー

``occonfig``パッケージはOpenCensusの初期化の手助けをします。
このパッケージを取り込んだアプリケーションは、環境変数やコマンドライン引数でOpenCensusのエクスポーターの設定ができるようになります。

## Instrucments

* 公開予定

## Exporters

* [opencensus-go-exporter-zap](https://github.com/future-architect/futureot/tree/master/exporters/opencensus-go-exporter-zap): [zap](https://godoc.org/go.uber.org/zap)経由でコンソールに出力するエクスポーターです。