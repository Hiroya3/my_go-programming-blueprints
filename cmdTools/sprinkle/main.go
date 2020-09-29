package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const otherWord = "*"

var transforms = []string{
	otherWord,
	otherWord,
	otherWord,
	otherWord,
	otherWord + "app",
	otherWord + "site",
	otherWord + "time",
	"get" + otherWord,
	"go" + otherWord,
	"lets " + otherWord,
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		t := transforms[rand.Intn(len(transforms))]
		fmt.Println(strings.Replace(t, otherWord, s.Text(), -1))
	}
}

//ファイルから読み込む用、*の処理をうまく行う必要あり
func readTransForms(fileName string) []string {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("ファイルの読み込みに失敗しました。\n エラ〜メッセージ：%s", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var result []string

	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	return result
}
