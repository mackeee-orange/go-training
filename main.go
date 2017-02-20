package main

import (
	"log"
	"net/http"
	"sync"
	"text/template"
	"path/filepath"
	"flag"
	"os"
	"github.com/makino18/training-go/playground_chat/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
)

type templateHandler struct {
	once	sync.Once
	filename string
	templ *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	addr := flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()

	gomniauth.SetSecurityKey("セキュリティキー")
	gomniauth.WithProviders(
		facebook.New("828855347266378", "e80426a214e562c683b096393c2a82b4", "http://localhost:8080/auth/callback/facebook"),
		github.New("88eb3e3ad78046e4b6b5", "c055fc391bd2a436360ff6c59b7cdbde7b1d2f1d", "http://localhost:8080/auth/callback/github"),
		google.New("314612542968-aa89bamhria784gdmtjuq0mvg3p84feo.apps.googleusercontent.com", "ymXsanZvYD4uklj6TZxzPS08", "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	//ルート
	//http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth", loginHandler)
	http.Handle("/room", r)

	go r.run()

	log.Println("Webサーバーを開始します。ポート：", *addr)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}


