package main

import (
	"GBP_Workthrough/chapter_1/trace"
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServerHTTP handles the HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse() // parse the flags

	// setup gomniauth
	gomniauth.SetSecurityKey("some long key")
	gomniauth.WithProviders(
		facebook.New("1653439198252194", "9e6e2b600cc9578f669892a24b132c2f", "http://localhost:8080/auth/callback/facebook"),
		github.New("043c7c4c40f9d7bf4778", "58f8c29b19761867661922ea1cb97bd43d5b00a1", "http://localhost:8080/auth/callback/github"),
		google.New("625130278991-q752b2rpfit512dl119tsfh4shq0f25h.apps.googleusercontent.com", "NxFsba2SFuyeFRuZGVxvXuC7", "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	// get the room going
	go r.run()

	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
