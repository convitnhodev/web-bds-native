package web

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bmizerany/pat"
	"github.com/deeincom/deeincom/internal/root"
	"github.com/golangcollege/sessions"
)

type router struct {
	*root.Cmd
	html    map[string]*template.Template
	session *sessions.Session
}

func init() {
	rand.Seed(time.Now().UnixNano())
	cmd := root.New("web")
	port := cmd.String("port", ":3000", "web listen port")
	cmd.Action(func() error {
		return run(cmd, *port)
	})
}

func (a *router) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := a.html[name]
	if !ok {
		a.render(w, r, "500.page.html", &templateData{})
		return
	}
	// apply global data, such as url, description etc..
	td.CurrentURL = r.RequestURI

	if r.Header.Get("X-Forwarded-For") == "" && (strings.Contains(r.Host, "local") || strings.Contains(r.Host, "127.0.0.1")) {
		td.Localhost = true
	}
	// get the current user id
	id := a.session.GetInt(r, "user")
	if id > 0 {
		// Get user by id
		u, err := a.App.Users.ID(id)
		if err != nil {
			log.Println(err)
			// user has been deleted? remove session anyway
			a.session.Remove(r, "user")
		} else {
			td.User = u
		}
	}

	buf := new(bytes.Buffer)
	if err := ts.Execute(buf, td); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	buf.WriteTo(w)
}

func (a *router) home(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "home.page.html", &templateData{
		HomePage: true,
	})
}

func (a *router) robots(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("./robots.txt")
	if err == nil {
		io.Copy(w, f)
		return
	}

	fmt.Fprintln(w, "User-agent: *")
	fmt.Fprintln(w, "Disallow: /")
}

func run(c *root.Cmd, port string) error {
	begin := time.Now().UnixNano()

	// runs all pending migrations before starting the app
	// failed to do so result in panic
	if err := c.App.Migration.Migrate(); err != nil {
		panic(err)
	}

	mux := pat.New()

	session := sessions.New([]byte("rat_la_bi_mat"))
	session.Lifetime = 24 * time.Hour * 30 // 1 month

	html, err := parseHTML(filepath.Join("ui", "html"))
	if err != nil {
		return err
	}

	a := &router{
		Cmd:     c,
		html:    html,
		session: session,
	}

	mux.Get("/", http.HandlerFunc(a.home))
	mux.Get("/robots.txt", http.HandlerFunc(a.robots))

	fs := http.FileServer(http.Dir(filepath.Join("ui", "static")))
	mux.Get("/static/", http.StripPrefix("/static", fs))

	srv := &http.Server{
		Addr:         port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  150 * time.Second,
		Handler:      session.Enable(mux),
	}

	end := time.Now().UnixNano()

	fmt.Printf("LOAD TIME: %f\n", float64(end-begin)/1000000000)
	return srv.ListenAndServe()
}
