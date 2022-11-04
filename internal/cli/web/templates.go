package web

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/deeincom/deeincom/config"
	"github.com/deeincom/deeincom/pkg/form"
	"github.com/deeincom/deeincom/pkg/models"
	"github.com/deeincom/deeincom/pkg/models/db"
	"github.com/microcosm-cc/bluemonday"
)

var strictPolicy = bluemonday.StrictPolicy()

type templateData struct {
	User       *models.User
	Users      []*models.User
	Pagination *db.Pagination

	Localhost  bool
	CurrentURL string

	Flash string

	Form     *form.Form
	Products []*models.Product
	Product  *models.Product

	Log  *models.Log
	Logs []*models.Log

	//Config
	Config *config.Config
}

var functions = template.FuncMap{
	"__":              translate,
	"upper":           strings.ToUpper,
	"lower":           strings.ToLower,
	"title":           strings.Title,
	"split":           strings.Split,
	"contains":        strings.Contains,
	"has_role":        hasRole,
	"html":            html,
	"buildPagination": buildPagination,
	"sureFind":        sureFind,
}

// sureFind always find an element in list l
// no matter what string s is
// it also find same element in `l` when passing same s
// return empty string if `l“ was empty
func sureFind(l []string, s string) string {
	if len(l) == 0 {
		return ""
	}

	r := []rune(s)
	total := 0
	for _, i := range r {
		total = total + int(i)
	}

	index := 0
	if total == 0 {
		index = rand.Intn(len(l) - 1)
	} else {
		index = total % len(l)
	}

	return l[index]
}

func html(s string) template.HTML {
	return template.HTML(s)
}

func buildPagination(url string, page int) string {
	if !strings.Contains(url, "?") {
		return fmt.Sprintf("%s?page=%d", url, page)
	}
	if !strings.Contains(url, "page") {
		return fmt.Sprintf("%s&page=%d", url, page)
	}
	if !strings.Contains(url, "&") {
		return strings.Split(url, "?")[0] + fmt.Sprintf("?page=%d", page)
	}
	s := []string{}
	if strings.Contains(url, "page") {
		cURL := strings.Split(url, "&")
		for _, v := range cURL {
			if strings.Contains(v, "page") {
				v = fmt.Sprintf("page=%d", page)
			}
			s = append(s, v)
		}
		return strings.Join(s, "&")
	}

	return ""
}

func hasRole(user *models.User, role string) bool {
	if user == nil {
		return false
	}

	for _, s := range user.Roles {
		if s == role {
			return true
		}
	}
	return false
}

func translate(s string) string {
	switch s {
	case "err_parse_form":
		return "Không thể lấy dữ liệu"

	case "err_invalid_form":
		return "Thông tin gởi lên không đúng"

	case "err_could_not_create_user":
		return "Đăng ký thất bại"

	case "err_could_not_verified_phone":
		return "Không thể xác thực điện thoại"

	case "err_could_not_verified_email":
		return "Không thể xác thực email"
	}
	return s
}

func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

func parseHTML(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "pages", "*.page.html"))
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		ts := template.New(name).Funcs(functions)

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "partials", "*.partial.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
