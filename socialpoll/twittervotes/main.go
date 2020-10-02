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
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}
	iter.Close()
	return options, iter.Err()
}

func main() {}
