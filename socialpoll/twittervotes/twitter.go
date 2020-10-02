package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/joeshaw/envdecode"
)

var conn net.Conn

//TwitterとのAPIの接続(net.Conn)を作成
func dial(netw, addr string) (net.Conn, error) {
	//はじめに閉じないと、異常終了等で接続が切れた場合、net.Connが彷徨う
	if conn != nil {
		conn.Close()
		conn = nil
	}

	netc, err := net.DialTimeout(netw, addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	conn = netc
	return netc, nil
}

var reader io.ReadCloser

//接続を切る
func closeConn() {
	if conn != nil {
		conn.Close()
	}
	if reader != nil {
		reader.Close()
	}
}

var (
	authClient *oauth.Client
	creds      *oauth.Credentials
)

//環境変数を読み込む
func setupTwitterAuth() {
	var ts struct {
		ConsumerKey    string `env:"SP_TWITTER_KEY",required`
		ConsumerSecret string `env:"SP_TWITTER_SECRET",required`
		AccessToken    string `env:"SP_TWITTER_ACCESSTOKEN",required`
		AccessSecret   string `env:"SP_TWITTER_ACCESSSECRET",required`
	}

	//環境変数読み込み
	if err := envdecode.Decode(&ts); err != nil {
		log.Fatalln(err)
	}

	creds = &oauth.Credentials{
		Token:  ts.AccessSecret,
		Secret: ts.AccessToken,
	}

	authClient = &oauth.Client{
		Credentials: oauth.Credentials{
			Token:  ts.ConsumerKey,
			Secret: ts.ConsumerSecret,
		},
	}
}

var (
	authSetupOnce sync.Once
	httpClient    *http.Client
)

func makeRequest(req *http.Request, params url.Values) (*http.Response, error) {
	//初期化コードは1回のみ
	authSetupOnce.Do(func() {
		setupTwitterAuth()
		httpClient = &http.Client{
			Transport: &http.Transport{
				//TODO：TCPコネクションの実現
				Dial: dial,
			},
		}
	})
	formEnc := params.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencode")
	req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))
	req.Header.Set("Authorization", authClient.AuthorizationHeader(creds, "POST", req.URL, params))
	return httpClient.Do(req)
}
