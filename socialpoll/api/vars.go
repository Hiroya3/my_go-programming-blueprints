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
