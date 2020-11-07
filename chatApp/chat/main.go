package main

import (
	"flag"
	"html/template"
	"log"
	"my_go-programming-blueprints/chatApp/trace"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

//現在のアクティブなAvatarの実装
var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UserGravatar,
}

//templateは１つのテンプレートを表します
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

//ServeHTTPはHTTPリクエストを処理します
//ServeHTTPメソッドがhttp.handlerインターフェイスに適合している＝内部的に呼び出されるもので、実装することで自分のsync.Onceが使えたりする
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() //フラグを解釈します

	//Gomniauthのセットアップ
	gomniauth.SetSecurityKey("t2wqrqRlrfAlkouK34")
	gomniauth.WithProviders(
		facebook.New("クライアントID", "秘密の値", "http://localhost:8080/auth/callback/facebook"),
		github.New("クライアントID", "秘密の値", "http://localhost:8080/auth/callback/github"),
		google.New("637491366567-mddrk8rpr6gc83cqqe1h9t9s9n64ikp4.apps.googleusercontent.com", "voUqFCVhXMSsHxOIaRk_dP17", "http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	//ルート(第2引数にはServeHTTPを持ったhandlerが必要)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars/"))))
	http.Handle("/room", r)
	//chatルームを開始します
	go r.run()
	//webサーバーを開始します
	log.Println("Webサーバーを起動します。ポート:", *addr)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
