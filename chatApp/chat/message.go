package main

import "time"

//messageは１つのメッセージを定義します
type message struct {
	Name      string
	Message   string
	When      time.Time
	AvatarURL string
}
