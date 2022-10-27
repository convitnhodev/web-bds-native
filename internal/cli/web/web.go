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
	"strconv"
	"strings"
	"time"

	"github.com/bmizerany/pat"
	"github.com/deeincom/deeincom/internal/cli/root"
	"github.com/deeincom/deeincom/pkg/form"
	"github.com/golangcollege/sessions"
)

var ui string   // đường dẫn đến thư mục ui
var port string // port web sẽ chạy

type router struct {
	*root.Cmd
	html    map[string]*template.Template
	session *sessions.Session
}

func init() {
	rand.Seed(time.Now().UnixNano())
	cmd := root.New("web")
	cmd.StringVar(&port, "port", ":3000", "port của web")
	cmd.StringVar(&ui, "ui", "ui/basic", "thư mục chứa theme cho ui")
	cmd.Action(func() error {
		return run(cmd)
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

func (a *router) productDetail(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "product.page.html", &templateData{})
}

func (a *router) verifyEmail(w http.ResponseWriter, r *http.Request) {
	f := form.New(r.PostForm)

	defer func() {
		a.render(w, r, "verify.email.page.html", &templateData{
			Form: f,
		})
	}()

	// nếu req là post thì đó là user muốn nhập lại email
	if r.Method == "POST" {
		f.Required("Email")

		if err := r.ParseForm(); err != nil {
			f.Errors.Add("err", "err_parse_form")
			return
		}

		if !f.Valid() {
			f.Errors.Add("err", "err_invalid_form")
			return
		}

		// thao tác này luôn thành công
		// hệ thống sẽ ko bao giờ báo cho user biết email có tồn tại hay không

		// lấy user bằng email
		user, err := a.App.Users.GetByEmail(f.Get("Email"))
		if err != nil {
			log.Println(err)
			return
		}

		// kiểm tra coi lần trước user yêu cầu verify email là lúc nào
		// mổi lần yêu cầu phải cách 60s
		if !user.SendVerifiedEmailAt.Add(60 * time.Second).UTC().Before(time.Now().UTC()) {
			return
		}
		// email.send()
	}

	// nếu ko có, hoặc ko parse dc iat, thì xem như expired
	iat, err := strconv.ParseInt(r.URL.Query().Get("iat"), 10, 64)
	if err != nil {
		log.Println(err)
		f.Errors.Add("err", "err_token_expired")
		return
	}

	// token của email sẽ expired sau 1 tuần
	isExpired := time.Unix(iat+3600*24*7, 0).UTC().Before(time.Now().UTC())
	if isExpired {
		log.Println(err)
		f.Errors.Add("err", "err_token_expired")
		return
	}

	token := r.URL.Query().Get("token")

	// tìm user với cái token này
	user, err := a.App.Users.GetByEmailToken(token)
	if err != nil {
		log.Println(err)
		f.Errors.Add("err", "err_token_invalid")
	}

	// update user đã verify email
	if err := a.App.Users.AddRole(user.ID, "verified_email"); err != nil {
		log.Println(err)
		f.Errors.Add("err", "err_could_not_verified_email")
	}
}

func (a *router) verifyPhone(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "verify.phone.page.html", &templateData{})
}

func (a *router) register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Println(err)
			a.render(w, r, "register.page.html", &templateData{})
			return
		}
	}
	a.render(w, r, "register.page.html", &templateData{})
}

func (a *router) login(w http.ResponseWriter, r *http.Request) {
	f := form.New(r.PostForm)
	f.Required("phone", "password")

	var ok bool
	defer func() {
		if !ok {
			a.render(w, r, "login.page.html", &templateData{
				Form: f,
			})
		}
	}()

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			f.Errors.Add("err", "err_parse_form")
			return
		}

		if !f.Valid() {
			f.Errors.Add("err", "err_invalid_form")
			return
		}

	}
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

func run(c *root.Cmd) error {
	begin := time.Now().UnixNano() // bắt đầu tính thời gian

	if err := c.App.Migration.Migrate(); err != nil {
		panic(err)
	}

	mux := pat.New()

	session := sessions.New([]byte("rat_la_bi_mat")) // uy tín
	session.Lifetime = 24 * time.Hour * 30           // user login 1 tháng mới bị out

	html, err := parseHTML(filepath.Join(ui, "html"))
	if err != nil {
		return err
	}

	a := &router{
		Cmd:     c,
		html:    html,
		session: session,
	}

	// homepage
	mux.Get("/", use(a.home))

	// product detail
	mux.Get("/real-estate/:slug", use(a.productDetail))

	// đăng ký
	mux.Post("/register", use(a.register))
	mux.Get("/register", use(a.register))

	// đăng nhập
	mux.Post("/login", use(a.login))
	mux.Get("/login", use(a.login))

	// verify
	mux.Get("/verify-email", use(a.verifyEmail))
	mux.Get("/verify-phone", use(a.verifyPhone))

	mux.Get("/robots.txt", use(a.robots))
	fs := http.FileServer(http.Dir(filepath.Join(ui, "static")))
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
