# my_go-programming-blueprints
[『Go言語によるWebアプリケーション開発』(Mat Ryer著、鵜飼文敏監訳、牧野聡訳; O'Reilly Japan)](https://github.com/oreilly-japan/go-programming-blueprints)の勉強用リポジトリ


## NSQをdockerで立てる
書籍ではローカルにnsqのインストールを行っていましたが、
自分はdockerで立てました。
(macを使っています)

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

### 3.nsqdを立てる
以下のコマンドでホストのIPアドレスを確認。

```
ifconfig | grep inet
```

([公式サイト](https://nsq.io/deployment/docker.html#run-nsqd)では `ifconfig | grep addr` でしたが自分は出ず)

IPアドレスを確認した上で以下のコマンドを実行。

```
docker run --name nsqd -p 4150:4150 -p 4151:4151 nsqio/nsq /nsqd --broadcast-address=<hostのipアドレス> --lookupd-tcp-address=<hostのipアドレス>:4160
```
