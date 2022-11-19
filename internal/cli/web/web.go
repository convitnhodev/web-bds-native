package web

import (
	"bytes"
	"crypto/md5"
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
	deein "github.com/deeincom/deeincom"
	"github.com/deeincom/deeincom/internal/cli/root"
	"github.com/deeincom/deeincom/pkg/email"
	"github.com/deeincom/deeincom/pkg/form"
	"github.com/deeincom/deeincom/pkg/helper"
	"github.com/deeincom/deeincom/pkg/models"
	"github.com/deeincom/deeincom/pkg/phone"
	"github.com/golangcollege/sessions"
)

var fe string   // đường dẫn đến thư mục theme cho fe
var be string   // đường dẫn đến thư mục theme cho be
var port string // port web sẽ chạy
var app deein.App

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

	attachments, err := a.App.Attachments.Find(product)
	if err != nil {
		log.Println(err)
		a.render(w, r, "404.page.html", &templateData{})
		return
	}

	comments, err := a.App.Comments.Slug(slug)
	if err != nil {
		comments = []*models.Comment{}
	}

	a.render(w, r, "detail.page.html", &templateData{
		Product:     product,
		Attachments: attachments,
		Comments:    comments,
	})
}

func (a *router) verifyEmailResult(w http.ResponseWriter, r *http.Request) {
	f := form.New(nil)

	defer func() {
		a.render(w, r, "result.email.page.html", &templateData{
			Form: f,
		})
	}()

	// check hash
	iat := r.URL.Query().Get("iat")
	token := r.URL.Query().Get("token")
	sign := r.URL.Query().Get("s")

	qs := fmt.Sprintf("deein:%s:%s:deein", token, iat)

	if fmt.Sprintf("%x", md5.Sum([]byte(qs))) != sign {
		log.Println("url sign không đúng")
		f.Errors.Add("err", "err_token_expired")
		return
	}

	// nếu ko có, hoặc ko parse dc iat, thì xem như expired
	issueAt, err := strconv.ParseInt(iat, 10, 64)
	if err != nil {
		log.Println(err)
		f.Errors.Add("err", "err_token_expired")
		return
	}

	// token của email sẽ expired sau 1 tuần
	isExpired := time.Unix(issueAt+3600*24*7, 0).UTC().Before(time.Now().UTC())
	if isExpired {
		log.Println(err)
		f.Errors.Add("err", "err_token_expired")
		return
	}

	// tìm user với cái token này
	user, err := a.App.Users.GetByEmailToken(token)
	if err != nil {
		log.Println(err)
		f.Errors.Add("err", "err_token_invalid")
	}

	// update user đã verify email
	if err := a.App.Users.AddRole(user, "verified_email"); err != nil {
		log.Println(err)
		f.Errors.Add("err", "err_could_not_verified_email")
	}

}

func (a *router) verifyPhoneResult(w http.ResponseWriter, r *http.Request) {
	f := form.New(nil)
	// f.Errors.Add("err", "err_parse_form")
	a.render(w, r, "result.phone.page.html", &templateData{
		Form: f,
	})
}

func (a *router) verifyEmail(w http.ResponseWriter, r *http.Request) {
	f := form.New(nil)

	ok := false
	defer func() {
		if !ok {
			a.render(w, r, "verify.email.page.html", &templateData{
				Form: f,
			})
		}
	}()

	// nếu user đã verify email rồi thì redirect về home
	if id := a.session.GetInt(r, "user"); id > 0 {
		if user, _ := a.App.Users.ID(fmt.Sprint(id)); hasRole(user, "verified_email") {
			ok = true
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}

	// nếu req là post thì đó là user muốn nhập lại email
	// dùng flash session để trả feedback, tránh bị nhầm với form
	if r.Method == "POST" {

		if err := r.ParseForm(); err != nil {
			f.Errors.Add("err", "err_parse_form")
			return
		}

		f.Values = r.PostForm
		f.Required("EmailToken")

		if !f.Valid() {
			log.Println("form invalid", f.Errors)
			return
		}

		// tìm user với cái token này
		user, err := a.App.Users.GetByEmailToken(f.Get("EmailToken"))
		if err != nil {
			log.Println(err)
			f.Errors.Add("err", "err_token_invalid")
		}

		// update user đã verify email
		if err := a.App.Users.AddRole(user, "verified_email"); err != nil {
			log.Println(err)
			f.Errors.Add("err", "err_could_not_verified_email")
		}

		ok = true
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}

}

func (a *router) verifyPhone(w http.ResponseWriter, r *http.Request) {
	f := form.New(nil)

	ok := false
	defer func() {
		if !ok {
			a.render(w, r, "verify.phone.page.html", &templateData{
				Form: f,
			})
		}
	}()

	// nếu user đã verify phone rồi thì redirect về home
	if id := a.session.GetInt(r, "user"); id > 0 {
		if user, _ := a.App.Users.ID(fmt.Sprint(id)); hasRole(user, "verified_phone") {
			ok = true
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			f.Errors.Add("err", "err_parse_form")
			return
		}

		f.Values = r.PostForm
		f.Required("PhoneToken")

		if !f.Valid() {
			log.Println("form invalid", f.Errors)
			return
		}

		// tìm user với cái token này
		user, err := a.App.Users.GetByPhoneToken(f.Get("PhoneToken"))
		if err != nil {
			log.Println(err)
			f.Errors.Add("err", "err_token_invalid")
		}

		// update user đã verify phone
		if err := a.App.Users.AddRole(user, "verified_phone"); err != nil {
			log.Println(err)
			f.Errors.Add("err", "err_could_not_verified_phone")
		}

		ok = true
		http.Redirect(w, r, "/verify/phone", http.StatusSeeOther)
	}
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

			// nếu user có nhập email thì gởi verify email luôn
			if f.Get("Email") != "" {
				// nhớ email
				user.Email = f.Get("Email")

				if err := a.App.Users.LogSendVerifyEmail(user); err != nil {
					log.Println(err)
				}

				// gởi email verify
				email.PostmarkApiToken = a.App.Config.PostmarkApiToken
				if err := email.SendVerifyEmail(user); err != nil {
					log.Println(err)
				}
			}
		}

		ok = true
		http.Redirect(w, r, "/verify-phone", http.StatusSeeOther)
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

func (a *router) forgotPassword(w http.ResponseWriter, r *http.Request) {
	f := form.New(nil)

	ok := false
	defer func() {
		if !ok {
			a.render(w, r, "forgot.password.page.html", &templateData{
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
		f.Required("Phone")

		if !f.Valid() {
			log.Println("form invalid", f.Errors)
			f.Errors.Add("err", "err_invalid_form")
			return
		}

		phoneNum := f.Get("Phone")
		user, err := a.App.Users.GetByPhone(phoneNum)

		if err != nil {
			f.Errors.Add("err", "err_invalid_form")
			return
		}

		if user != nil &&
			(user.RPTExpiredAt == nil || user.RPTExpiredAt.Before(time.Now())) {
			token := helper.RandDigitString(5)
			err := a.App.Users.ResetPasswordByPhone(phoneNum, token)

			if err != nil {
				log.Println("Reset fail", err)
				f.Errors.Add("err", "err_invalid_form")
				return
			}

			if err := phone.SendResetPwdPhone(phoneNum, token); err != nil {
				log.Println(err)
			}

		}

		ok = true
		a.session.Put(r, "reset_pwd_phone", phoneNum)
		http.Redirect(w, r, "/reset-password", http.StatusSeeOther)
	}
}

func (a *router) resetPassword(w http.ResponseWriter, r *http.Request) {
	f := form.New(nil)
	ok := false

	// Nếu không chứa reseting phone thì không show form
	phoneNum := a.session.GetString(r, "reset_pwd_phone")
	userId := a.session.GetInt(r, "user")

	if phoneNum == "" && userId <= 0 {
		ok = true
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	defer func() {
		if !ok {
			a.render(w, r, "reset.password.page.html", &templateData{
				Form:              f,
				IsResetPwdByToken: phoneNum != "",
			})
		}
	}()

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			f.Errors.Add("err", "err_parse_form")
			return
		}

		f.Values = r.PostForm
		f.Required("Token", "Password", "RepeatPassword")

		if f.Get("Password") != f.Get("RepeatPassword") {
			f.Errors.Add("err", "err_invalid_repeatpassword")
			return
		}

		if !f.Valid() {
			log.Println("form invalid", f.Errors)
			f.Errors.Add("err", "err_invalid_form")
			return
		}

		if phoneNum != "" {
			// Trường hợp reset by phone
			user, _ := a.App.Users.GetByPhone(phoneNum)

			if user != nil {
				err := a.App.Users.CompareHashAndPassword(user.ResetPasswordToken, f.Get("Token"))

				if err != nil || user.RPTExpiredAt.Before(time.Now()) {
					log.Println("Token invalid", user.ResetPasswordToken, user.RPTExpiredAt)
					f.Errors.Add("err", "err_invalid_token")
					return
				}

				a.App.Users.UpdateNewPassword(fmt.Sprint(user.ID), f.Get("Password"))
			}
		} else if userId > 0 {
			// Trường hợp reset by old password
			user, _ := a.App.Users.ID(fmt.Sprint(userId))

			if user != nil {
				err := a.App.Users.CompareHashAndPassword(user.Password, f.Get("Token"))
				if err != nil {
					log.Println("Token invalid", f.Get("Token"))
					f.Errors.Add("err", "err_invalid_token")
					return
				}

				a.App.Users.UpdateNewPassword(fmt.Sprint(user.ID), f.Get("Password"))
			}
		}

		a.session.Pop(r, "reset_pwd_phone")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func (a *router) islogined(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := a.session.GetInt(r, "user")
		if id == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *router) createComment(w http.ResponseWriter, r *http.Request) {
	f := form.New(nil)
	slug := r.URL.Query().Get(":slug")
	id := a.session.GetInt(r, "user")

	ok := false
	defer func() {
		if !ok {
			http.Error(w, "bad request", 400)
		}
	}()

	if err := r.ParseForm(); err != nil {
		f.Errors.Add("err", "err_parse_form")
		return
	}

	f.Values = r.PostForm
	f.Set("Slug", slug)
	f.Set("UserId", fmt.Sprint(id))
	f.Required("Message")

	if !f.Valid() {
		log.Println("form invalid", f.Errors)
		return
	}

	if _, err := a.App.Comments.Create(f); err != nil {
		log.Println(err)
		f.Errors.Add("err", "could_not_update_comment")
		return
	}

	ok = true

	w.Write([]byte("Ok"))
}

func (a *router) uploadKYC(w http.ResponseWriter, r *http.Request) {
	f := form.New(nil)
	userId := a.session.GetInt(r, "user")
	user, err := a.App.Users.ID(fmt.Sprint(userId))
	ok := false

	if err != nil {
		a.render(w, r, "500.page.html", &templateData{})
		return
	}

	// TODO: check user có verified phone chưa
	if !hasRole(user, "verify_phone") {
		http.Redirect(w, r, "/upgrade-user?to=verified_id", http.StatusSeeOther)
		return
	}

	defer func() {
		if !ok {
			a.render(w, r, "kyc.page.html", &templateData{
				Form: f,
				User: user,
			})
		}
	}()

	if r.Method == "POST" {
		r.Body = http.MaxBytesReader(w, r.Body, 30<<20)
		if err := r.ParseMultipartForm(30 << 20); err != nil {
			log.Println(err)
			f.Errors.Add("err", "err_parse_form")
			return
		}

		f.Values = r.PostForm
		f.Required("FrontIdentityCard")
		f.Required("BackIdentityCard")
		f.Required("SelfieImage")

		if !f.Valid() {
			log.Println("form invalid", f.Errors)
			return
		}

		// Front file
		frontFile, handler, err := r.FormFile("FrontIdentityCard")
		if err != nil && http.ErrMissingFile != err {
			log.Println(err)
			f.Errors.Add("err", "err_could_not_upload")
			return
		}
		defer frontFile.Close()

		frontFileName, err := a.App.LocalFile.UploadFile(fmt.Sprintf("users.%d/", userId), frontFile, handler)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f.Set("FrontIdentityCard", *frontFileName)

		// Back file
		backFile, handler, err := r.FormFile("FrontIdentityCard")
		if err != nil && http.ErrMissingFile != err {
			log.Println(err)
			f.Errors.Add("err", "err_could_not_upload")
			return
		}
		defer frontFile.Close()

		backFileName, err := a.App.LocalFile.UploadFile(fmt.Sprintf("users.%d/", userId), backFile, handler)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f.Set("BackIdentityCard", *backFileName)

		// Selfie file
		selfieFile, handler, err := r.FormFile("FrontIdentityCard")
		if err != nil && http.ErrMissingFile != err {
			log.Println(err)
			f.Errors.Add("err", "err_could_not_upload")
			return
		}
		defer frontFile.Close()

		selfieFileName, err := a.App.LocalFile.UploadFile(fmt.Sprintf("users.%d/", userId), selfieFile, handler)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f.Set("SelfieImage", *selfieFileName)

		if err := a.App.Users.UpdateKYCStatus(fmt.Sprint(userId), "submited_kyc"); err != nil {
			log.Println(err)
			f.Errors.Add("err", "could_not_kyc")
			return
		}

		if err := a.App.KYC.SubmitKYC(fmt.Sprint(userId), f); err != nil {
			log.Println(err)
			f.Errors.Add("err", "could_not_kyc")
			return
		}

		ok = true
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (a *router) upgradeUser(w http.ResponseWriter, r *http.Request) {
	to := r.URL.Query().Get("to")
	userId := a.session.GetInt(r, "user")
	user, err := a.App.Users.ID(fmt.Sprint(userId))

	log.Println(userId)

	if err != nil {
		a.render(w, r, "500.page.html", &templateData{})
		return
	}

	goLink := "/"
	switch to {
	case "verified_id":
		if hasRole(user, "verify_phone") {
			if !hasRole(user, "verified_id") {
				goLink = "/kyc"
			}
		} else {
			goLink = "/verify-phone"
		}
	case "deein_partner":
		if hasRole(user, "verified_id") {
			if !hasRole(user, "deein_partner") {
				goLink = "/apply-partner"
			}
		} else {
			goLink = "/kyc"
		}
	}

	http.Redirect(w, r, goLink, http.StatusSeeOther)
}

func (a *router) applyPartner(w http.ResponseWriter, r *http.Request) {
	f := form.New(nil)
	userId := a.session.GetInt(r, "user")
	user, err := a.App.Users.ID(fmt.Sprint(userId))
	ok := false

	if err != nil {
		a.render(w, r, "500.page.html", &templateData{})
		return
	}

	if !hasRole(user, "verified_id") {
		http.Redirect(w, r, "/upgrade-user?to=deein_partner", http.StatusSeeOther)
		return
	}

	defer func() {
		if !ok {
			a.render(w, r, "apply.partner.page.html", &templateData{
				Form: f,
				User: user,
			})
		}
	}()

	if r.Method == "POST" {
		r.Body = http.MaxBytesReader(w, r.Body, 30<<20)
		if err := r.ParseMultipartForm(30 << 20); err != nil {
			log.Println(err)
			f.Errors.Add("err", "err_parse_form")
			return
		}

		f.Values = r.PostForm
		f.Required("Message")

		if !f.Valid() {
			log.Println("form invalid", f.Errors)
			return
		}

		// Front file
		file, handler, err := r.FormFile("UploadFile")
		if err != nil && http.ErrMissingFile != err {
			log.Println(err)
			f.Errors.Add("err", "err_could_not_upload")
			return
		}
		defer file.Close()

		cvFileName, err := a.App.LocalFile.UploadFile(fmt.Sprintf("partner.%d/", userId), file, handler)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f.Set("CVLink", *cvFileName)

		if err := a.App.Users.UpdatePartnerStatus(fmt.Sprint(userId), "apply"); err != nil {
			log.Println(err)
			f.Errors.Add("err", "could_not_submit_partner")
			return
		}

		if err := a.App.Partner.SubmitPartner(fmt.Sprint(userId), f); err != nil {
			log.Println(err)
			f.Errors.Add("err", "could_not_submit_partner")
			return
		}

		ok = true
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

	app = *a.App

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

	// Quên pass
	mux.Get("/forgot-password", use(a.forgotPassword))
	mux.Post("/forgot-password", use(a.forgotPassword))

	mux.Get("/reset-password", use(a.resetPassword))
	mux.Post("/reset-password", use(a.resetPassword))

	// out
	mux.Get("/logout", use(a.logout))

	// verify
	mux.Post("/verify-email", use(a.verifyEmail))
	mux.Get("/verify-email", use(a.verifyEmail))
	mux.Get("/verify/email", use(a.verifyEmailResult))

	mux.Post("/verify-phone", use(a.verifyPhone))
	mux.Get("/verify-phone", use(a.verifyPhone))
	mux.Get("/verify/phone", use(a.verifyPhoneResult))

	// comment
	mux.Post("/comments/:slug/create", use(a.createComment, a.islogined))

	// upgrade
	mux.Get("/upgrade-user", use(a.upgradeUser, a.islogined))
	mux.Get("/apply-partner", use(a.applyPartner, a.islogined))
	mux.Post("/apply-partner", use(a.applyPartner, a.islogined))

	// kyc
	mux.Get("/kyc", use(a.uploadKYC, a.islogined))
	mux.Post("/kyc", use(a.uploadKYC, a.islogined))

	mux.Get("/robots.txt", use(a.robots))

	// mấy trang linh tinh legal
	mux.Get("/privacy-notice", use(a.privacy))
	mux.Get("/terms-of-service", use(a.terms))

	fefs := http.FileServer(http.Dir(filepath.Join(fe, "static")))
	befs := http.FileServer(http.Dir(filepath.Join(be, "static")))
	mux.Get("/static/", http.StripPrefix("/static", fefs))
	mux.Get("/be/", http.StripPrefix("/be", befs))

	registerAdminRoute(mux, a)
	registerAjaxRoute(mux, a)

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
