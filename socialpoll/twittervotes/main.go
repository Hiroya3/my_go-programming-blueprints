package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bitly/go-nsq"
	"gopkg.in/mgo.v2"
)

var db *mgo.Session

//mongoDBへの接続
func dialdb() error {
	var err error
	log.Println("MongiDBにダイヤル中：localhost")
	db, err = mgo.Dial("localhost")
	return err
}

//mongoDBへの接続を解除
func closedb() {
	db.Clone()
	log.Println("データベース接続が閉じられました")
}

type poll struct {
	Options []string
}

func loadOptions() ([]string, error) {
	var options []string
	//ballotsデータベースに含まれるpollコレクション
	//find(nil)はフィルタリングを行わない
	iter := db.DB("ballots").C("polls").Find(nil).Iter()
	var p poll
	//nextで次にあった場合には引数に結果を入れる(スライスだが、1つづつしか入らない(?))
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}
	iter.Close()
	return options, iter.Err()
}

func publishVotes(votes <-chan string) <-chan struct{} {
	stopchan := make(chan struct{}, 1)
	pub, _ := nsq.NewProducer("localhost:4150", nsq.NewConfig())
	go func() {
		for vote := range votes {
			//投票内容（mongoDBで指定した言葉が含まれるtweet）をパブリッシュします
			pub.Publish("votes", []byte(vote))
		}
		log.Fatalln("Publisher : 停止中です")
		pub.Stop()
		log.Fatalln("Publisher : 停止しました")
		//deferでも良い
		stopchan <- struct{}{}
	}()
	return stopchan
}

//概要
//mongoDBにある言葉を検索し、ランキングにする
//本書では投票とあるのは、mongoDBにあらかじめある言葉を検索し数を数えることが
//その言葉に対する投票に見立てている
func main() {
	var stoplock sync.Mutex
	stop := false
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		stoplock.Lock()
		stop = true
		stoplock.Unlock()
		log.Println("停止します...")
		stopChan <- struct{}{}
		closeConn()
	}()

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	if err := dialdb(); err != nil {
		log.Fatalln("MongoDB屁のダイヤルに失敗しました：", err)
	}
	defer closedb()

	//処理を開始します

	//投票結果のためのチャネル
	votes := make(chan string)
	publisherStoppedChan := publishVotes(votes)
	twitterStoppedChan := startTwitterStream(stopChan, votes)
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			closeConn()
			stoplock.Lock()
			if stop {
				stoplock.Unlock()
				break
			}
			stoplock.Unlock()
		}
	}()
	<-twitterStoppedChan
	close(votes)
	<-publisherStoppedChan
}
