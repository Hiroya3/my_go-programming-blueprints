# 前提環境
- mac
- Dockerはインストール済み

# DockerでNSQとMongoDBコンテナを立てる
## NSQをdockerで立てる
書籍(5章から)ではローカルにnsqのインストールを行っていましたが、
自分はdockerで立てました。

### 立て方
基本的にNSQの[公式サイト](https://nsq.io/deployment/docker.html)の通りに行いました。
`nsqlookupd` と `nsqd` をrunします。
#### 1.イメージをpullする
以下のコマンドで `nsqio/nsq` をpullする。

```
docker pull nsqio/nsq
```

#### 2.nsqlookupdを立てる
以下のコマンドで `nsqlookupd` を立てる。

```
docker run --name lookupd -p 4160:4160 -p 4161:4161 nsqio/nsq /nsqlookupd
```

以下の出力が出る。

```
[nsqlookupd] 2020/09/29 13:02:11.974364 INFO: nsqlookupd v1.2.0 (built w/go1.12.9)
[nsqlookupd] 2020/09/29 13:02:11.975027 INFO: HTTP: listening on [::]:4161
[nsqlookupd] 2020/09/29 13:02:11.975103 INFO: TCP: listening on [::]:4160
```

#### 3.nsqdを立てる
以下のコマンドでホストのIPアドレスを確認。

```
ifconfig | grep inet
```

([公式サイト](https://nsq.io/deployment/docker.html#run-nsqd)では `ifconfig | grep addr` でしたが自分は出ず)

IPアドレスを確認した上で以下のコマンドを実行。

```
docker run --name nsqd -p 4150:4150 -p 4151:4151 nsqio/nsq /nsqd --broadcast-address=<hostのipアドレス> --lookupd-tcp-address=<hostのipアドレス>:4160
```

## MongoDBをdockerで立てる
書籍ではローカルにMongoDBのインストールを行っていましたが、
自分はdockerで立てました。
(macを使っています)

### 立て方
基本的に[公式サイト](https://hub.docker.com/_/mongo)の通りに行いました。

### 1.イメージをpullする
以下のコマンドで `mongo` をpullする。

```
docker pull mongo
```

### 2.mongoを立てる
以下のコマンドで `mongo` を立てます。

```
docker run --name some-mongo -p27017:27017 -d mongo
```

`some-mongo` はコンテナの名前。

# MongoDBにデータをinsert
上記で立てたMongoDBに対してデータをinsertする方法。

## 1.dockerコンテナの起動
MongoDBが起動されていない場合は起動します。
起動済みの場合は2に進んでください。
下記コマンドでdockerIDを確認します。

```
docker ps -a
```

`CONTAINER ID` を確認してください。

その後、以下のコマンドを実行してdockerコンテナを起動します。

```
docker start {CONTAINER ID}
```

## 2.MongoDBにアクセス
[公式リファレンス](https://hub.docker.com/_/mongo)を参照しています。

以下のコマンドを実行し、dockerコンテナのシェルにアクセスする。

```
docker exec -it {CONTAINER ID} bash
```

## 3.データベースの作成
MongoDBのシェル内で以下のコマンドを実行。

```
use {DB名}
```

## 4.データのインサート
下記のコマンドで作成したDBにデータのコレクションにデータを追加します。
※作成したDBに切り替わっていることを確認してください。

```
db.{コレクション名}.insert({データ})
```
