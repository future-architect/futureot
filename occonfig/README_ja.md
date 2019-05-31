# occonfig: OpenCensus Configurator

このパッケージは次のようなOpenCensusの初期化の手段を提供します

* 環境変数
* JSONファイル
* コマンドラインオプション

## occonfigが組み込まれたアプリケーションのユーザー向け説明

### 環境変数経由

* ``OC_SERVICE_NAME``: 出力先のレポートに出力されるサービス名。デフォルトはコマンドラインのプログラム名。

* ``OC_SERVICE_URL``:  サービスのURL。いくつかのトレーサーはURLで結果をフィルタリングする。

* ``OC_CONFIG_JSON``: JSON設定ファイルのパス（後述）

* ``OC_TRACE_EXPORTER``: トレーシングに必要

   * ``stackdriver://demo-project-id``: Stackdriver
   * ``sd://demo-project-id`` : Stackdriverの短縮系
   * ``datadog://localhost:8126`` もしくは ``dd://localhost:8126`` : DataDog
   * ``datadog`` もしくは ``dd`` : DataDog (デフォルトのホスト:ポートはlocalhost:8126)
   * ``xray``: AWS X-Ray
   * ``jaeger://localhost:6831`` : Jaeger
   * ``jaeger://localhost`` : Jaeger (デフォルトのポートは6831)
   * ``jaeger`` : Jaeger (デフォルトのホスト:ポートはlocalhost:6831)
   * ``zipkin://localhost:9411/api/v2/spans`` : Zipkin
   * ``zipkin://localhost/api/v2/spans`` : Zipkin (デフォルトポートは9411)
   * ``zipkin://localhost`` : Zipkin (デフォルトポートは9411, デフォルトパスは/api/v2/spans)
   * ``zap``: [zap](https://godoc.org/go.uber.org/zap)経由でコンソールに出力
   * ``honeycomb`` : HoneyComb

* ``OC_TRACE_SAMPLER``

   * ``always``: デフォルト値
   * ``never``: 出力しない
   * 浮動小数点数 (0-1): 確率的なサンプラー

* ``OC_HONEYCOMB_WRITE_KEY``: honeycomb.io APIキー。もし、値が　``file://``　から始まっていたら、ローカルのファイルを探索する。

### 一般的な利用方法

#### 環境変数経由

```bash
$ export OC_TRACE_EXPORTER=stackdriver://demo-project-id
$ export OC_SERVICE_NAME=my-service
$ ./your-program
```

#### flagパッケージのコマンドラインフラグ経由

* 共通設定

   * ``-oc-service-name``: サービス名
   * ``-oc-service-url``: サービスURL
   * ``-oc-config-json``: JSON形式の設定ファイルのパス（後述）

* トレースの設定

   * ``-oc-honeycomb-write-key``: honeycomb.ioのキーファイルパス
   * ``-oc-trace-exporter``: エクスポーター設定

```bash
# flagパッケージサポート
$ ./your-program -oc-trace-exporter stackdriver://demo-project-id -oc-service-name my-service
```

#### kingpin.v2パッケージのコマンドラインフラグ経由

* 共通設定

   * ``--oc-service-name``: サービス名
   * ``--oc-service-url``: サービスURL
   * ``--oc-config-json``: JSON形式の設定ファイルのパス（後述）

* トレースの設定

   * ``--oc-trace-exporter``: エクスポーターの設定
   * ``--oc-honeycomb-write-key``: honeycomb.ioのキーファイルパス

```
# kingpin.v2パッケージサポート
$ ./your-program --oc-trace-exporter=xray --oc-service-name=my-service
```

### JSONファイルフォーマット

設定ファイルのパスは``-oc-config-json`` (flagパッケージ利用時),  ``--oc-config-json`` (kingpin.v2パッケージ利用時)のオプションで指定できます。

extendsで、ベースとなるJSONを設定できます。

```json
{
  "service-name": "my-awesome-service",
  "service-url":  "http://localhost:8080",
  "extends": "../config.json",
  "trace": {
    "exporter": "stackdriver://demo-project-id",
    "honeycomb-write-key": "honeycomb.key",
    "sampling": "always"
  }
}
```

## 設定の優先度(小さい数字が優先度高)

1. コマンドラインオプション
2. コマンドラインオプションの``--oc-config-json``で指定されたJSONファイル
3. コマンドラインオプションの``--oc-config-json``で指定されたJSONファイルの``extends``で指定されたファイル
4. 環境変数
5. ``OC_CONFIG_JSON``の環境変数で指定されたJSONファイル
6. ``OC_CONFIG_JSON``の環境変数で指定されたJSONファイルの``extends``で指定されたファイル

## ローカル開発時のツール設定

### Jaegerの一般的な設定方法

#### Docker

```bash
$ docker run -d --name jaeger -p 14268:14268 -p 16686:16686 jaegertracing/all-in-one:1.12
$ OC_TRACE_EXPORTER=jaeger ./your-program
```

#### docker-compose

```yaml
version: '3'
services:
  jaeger:
    image: jaegertracing/all-in-one:1.12
    ports:
      - 16686:16686  # Webコンソール用
      - 14268:14268  # もし、開発中のサービスがdocker-composeの外で稼働している場合に必要
  your-service:
    image: your-service
    ports:
      - 8080:8080
    environment:
      - OC_TRACE_EXPORTER=jaeger://jaeger
      - OC_SERVICE_URL=http://localhost:8080
    depends_on:
      - jaeger
```

### Zipkinの一般的な設定方法

#### Docker

```bash
$ docker run -d -p 9411:9411 openzipkin/zipkin
$ OC_TRACE_EXPORTER=zipkin ./your-program
```

#### docker-compose

```yaml
version: '3'
services:
  zipkin:
    image: openzipkin/zipkin
    ports:
      - 9411:9411
  your-service:
    image: your-service
    ports:
      - 8080:8080
    environment:
      - OC_TRACE_EXPORTER=zipkin://zipkin
      - OC_SERVICE_URL=http://localhost:8080
    depends_on:
      - zipkin
```

## occonfigライブラリ利用者の設定方法

```go
package main

import (
	"flag"
	"github.com/shibukawa/occonfig"
	"gopkg.in/alecthomas/kingpin.v2"
)

// 使用方法1: 環境変数のみを利用
func main() {
	// 環境変数の情報を取得して初期化
	finalizer, err := occonfig.Init(occonfig.Trace | occonfig.Stats)
	if err != nil {
		panic(err)
	}
	defer finalizer.Close()
	
	// アプリケーションコードはここから
}

// 使用方法2: 環境変数とflagパッケージを利用
func main() {
	// flag.Parse()を呼び出す前に次の関数を実行
	occonfig.UseFlag(occonfig.Trace)
	
	flag.Parse()
	
	// その後Init()を呼び出す
	finalizer, err := occonfig.Init(occonfig.Trace)
	if err != nil {
		panic(err)
	}
	defer finalizer.Close()

	// アプリケーションコードはここから
}

// 使用方法2: 環境変数とkingpin.v2パッケージを利用
func main() {
	// kingpin.Parse()を呼び出す前に次の関数を実行
	occonfig.UseKingpin(occonfig.Stats)
	
	kingpin.Parse()
	
	// その後Init()を呼び出す
	finalizer, err := occonfig.Init(occonfig.Stats)
	if err != nil {
		panic(err)
	}
	defer finalizer.Close()

	// アプリケーションコードはここから
}
```