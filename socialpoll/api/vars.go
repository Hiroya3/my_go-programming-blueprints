package main

import (
	"net/http"
	"sync"
)

var (
	//複数のリクエストが同時にアクセスするためlock機構が必要
	varsLock sync.RWMutex
	//key:*http.Request  value :別のmap(map[string]interface{})
	vars map[*http.Request]map[string]interface{}
)

func OpenVars(r *http.Request) {
	varsLock.Lock()
	if vars == nil {
		vars = map[*http.Request]map[string]interface{}{}
	}
	vars[r] = map[string]interface{}{}
	varsLock.Unlock()
}

//メモリリークを防ぐためメモリの解放を行う
func CloseVars(r *http.Request) {
	varsLock.Lock()
	delete(vars, r)
	varsLock.Unlock()
}
