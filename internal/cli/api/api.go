package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/deeincom/deeincom/internal/cli/root"
)

func init() {
	cmd := root.New("api")
	port := cmd.String("port", ":4000", "api listen port")
	cmd.Action(func() error {
		return run(cmd, *port)
	})
}

type router struct {
	*root.Cmd
}

func (c *router) error(w http.ResponseWriter, s string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)

	result := map[string]interface{}{
		"_status":     code,
		"_success":    false,
		"_message":    s,
		"_data":       []struct{}{},
		"_pagination": struct{}{},
	}

	b, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		http.Error(w, "json decode err", 500)
		return
	}
	w.Write(b)
}

func (a *router) json(w http.ResponseWriter, v interface{}, e interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	result := map[string]interface{}{
		"_status":     200,
		"_success":    true,
		"_message":    "",
		"_data":       v,
		"_pagination": e,
	}

	b, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		http.Error(w, "json decode err", 500)
		return
	}
	w.Write(b)
}

// api should disallow all bot
func (a *router) robots(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User-agent: *")
	fmt.Fprintln(w, "Disallow: /")
}

func (a *router) home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("nothing here !"))
}

func run(c *root.Cmd, port string) error {
	// runs all pending migrations before starting the app
	// failed to do so result in panic
	if err := c.App.Migration.Migrate(); err != nil {
		panic(err)
	}

	a := &router{c}
	mux := http.NewServeMux()

	mux.HandleFunc("/", a.home)
	mux.HandleFunc("/robots.txt", a.robots)

	srv := &http.Server{
		Addr:         port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  150 * time.Second,
		Handler:      mux,
	}

	return srv.ListenAndServe()
}
