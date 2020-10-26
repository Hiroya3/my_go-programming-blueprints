package main

import (
	"encoding/json"
	"my_go-programming-blueprints/meander"
	"net/http"
	"runtime"
)

func main() {
	//プログラムから利用できるCPU数の最大値を指定。今回は全てのCPUを利用するように指定(go1.4までは必要だった)
	runtime.GOMAXPROCS(runtime.NumCPU())
	//meander.APIKey = "TODO"
	http.HandleFunc("/journeys", func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, meander.Journeys)
	})
	http.ListenAndServe(":8080", http.DefaultServeMux)
}

//抽象化
func respond(w http.ResponseWriter, r *http.Request, data []interface{}) error {
	publicData := make([]interface{}, len(data))
	for i, d := range data {
		//public.goのPublic()を使い、Publicが実装されているか確認
		publicData[i] = meander.Public(d)
	}
	return json.NewEncoder(w).Encode(publicData)
}