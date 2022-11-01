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
	"github.com/deeincom/deeincom/pkg/email"
	"github.com/deeincom/deeincom/pkg/form"
	"github.com/golangcollege/sessions"
)

var fe string   // đường dẫn đến thư mục theme cho fe
var be string   // đường dẫn đến thư mục theme cho be
var port string // port web sẽ chạy

type router struct {
	*root.Cmd
	fe      map[string]*template.Template
	be      map[string]*template.Template
	session *sessions.Session
}

func init() {
	rand.Seed(time.Now().UnixNano())
	cmd := root.New("web")
	cmd.StringVar(&port, "port", ":3000", "port của web")
	cmd.StringVar(&fe, "fe", "ui/basic", "thư mục chứa theme cho fe")
	cmd.StringVar(&be, "be", "ui/admin", "thư mục chứa theme cho be")
	cmd.Action(func() error {
		return run(cmd)
	})
}

func (a *router) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := a.fe[name]
	if !ok {
		a.render(w, r, "500.page.html", &templateData{})
		return
	}
	// apply global data, such as url, description etc..
	td.CurrentURL = r.RequestURI
	td.Config = a.App.Config

	// flash msg
	td.Flash = a.session.PopString(r, "flash")

	if r.Header.Get("X-Forwarded-For") == "" && (strings.Contains(r.Host, "local") || strings.Contains(r.Host, "127.0.0.1")) {
		td.Localhost = true
	}
	// get the current user id
	id := a.session.GetInt(r, "user")
	if id > 0 {
		// Get user by id
		u, err := a.App.Users.ID(fmt.Sprint(id))
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
func (a *router) terms(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "terms.page.html", &templateData{})
}
func (a *router) privacy(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "privacy.page.html", &templateData{})
}

func (a *router) home(w http.ResponseWriter, r *http.Request) {
	products, err := a.App.Products.Find()
	if err != nil {
		log.Println(err)
		a.render(w, r, "404.page.html", &templateData{})
		return
	}

	a.render(w, r, "home.page.html", &templateData{
		Products: products,
	})
}

func (a *router) productDetail(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get(":slug")
	product, err := a.App.Products.GetBySlug(slug)
	if err != nil {
		log.Println(err)
		a.render(w, r, "404.page.html", &templateData{})
		return
	}

	a.render(w, r, "detail.page.html", &templateData{
		Product: product,
	})
}

func (a *router) verifyEmail(w http.ResponseWriter, r *http.Request) {
	f := form.New(r.PostForm)

	defer func() {
		a.render(w, r, "verify.email.page.html", &templateData{
			Form: f,
		})
	}()

	// nếu req là post thì đó là user muốn nhập lại email
	// dùng flash session để trả feedback, tránh bị nhầm với form
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
		a.session.Put(r, "flash", "msg_verify_email_success")

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

		// bắn email
		email.Send()
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

	// tìm user với cái token này
	token := r.URL.Query().Get("token")
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
	f := form.New(nil)

	var ok bool
	defer func() {
		if !ok {
			a.render(w, r, "register.page.html", &templateData{
				Form: f,
			})
			return
		}

	}()

	// nếu đã login thì ko show nữa
	if a.session.GetInt(r, "user") > 0 {
		ok = true
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Println(err)
			f.Errors.Add("err", "err_parse_form")
			return
		}

		f.Values = r.PostForm
		f.Required("Phone", "Password", "FirstName", "LastName")

		if !f.Valid() {
			log.Println("form invalid", f.Errors)
			f.Errors.Add("err", "err_invalid_form")
			return
		}

		if user, err := a.App.Users.Create(f); err != nil {
			log.Println(err)
			f.Errors.Add("err", "err_could_not_create_user")
			return
		} else {
			a.session.Put(r, "user", user.ID)
		}

		ok = true
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	}
}

func (a *router) logout(w http.ResponseWriter, r *http.Request) {
	a.session.Remove(r, "user")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *router) login(w http.ResponseWriter, r *http.Request) {
	f := form.New(nil)

	var ok bool
	defer func() {
		if !ok {
			a.render(w, r, "login.page.html", &templateData{
				Form: f,
			})
		}
	}()

	// nếu đã login thì ko show nữa
	if a.session.GetInt(r, "user") > 0 {
		ok = true
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			f.Errors.Add("err", "err_parse_form")
			return
		}

		f.Values = r.PostForm
		f.Required("Phone", "Password")

		if !f.Valid() {
			log.Println("form invalid", f.Errors)
			f.Errors.Add("err", "err_invalid_form")
			return
		}

		user, err := a.App.Users.Auth(f)
		if err != nil {
			log.Println(err)
			f.Errors.Add("err", "msg_invalid_login")
			return
		}
		fmt.Println("DEBUG", "login", "success", user.ID)
		ok = true
		a.session.Put(r, "user", user.ID)
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
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

	feHTML, err := parseHTML(filepath.Join(fe, "html"))
	if err != nil {
		return err
	}
	beHTML, err := parseHTML(filepath.Join(be, "html"))
	if err != nil {
		return err
	}

	a := &router{
		Cmd:     c,
		fe:      feHTML,
		be:      beHTML,
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

	// out
	mux.Get("/logout", use(a.logout))

	// verify
	mux.Get("/verify-email", use(a.verifyEmail))
	mux.Get("/verify-phone", use(a.verifyPhone))

	mux.Get("/robots.txt", use(a.robots))

	// mấy trang linh tinh legal
	mux.Get("/privacy-notice", use(a.privacy))
	mux.Get("/terms-of-service", use(a.terms))

	fefs := http.FileServer(http.Dir(filepath.Join(fe, "static")))
	befs := http.FileServer(http.Dir(filepath.Join(be, "static")))
	mux.Get("/static/", http.StripPrefix("/static", fefs))
	mux.Get("/be/", http.StripPrefix("/be", befs))

	registerAdminRoute(mux, a)

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
