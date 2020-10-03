package main

import (
	"log"

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

//概要
//mongoDBにある言葉を検索し、ランキングにする
//本書では投票とあるのは、mongoDBにあらかじめある言葉を検索し数を数えることが
//その言葉に対する投票に見立てている
func main() {}
