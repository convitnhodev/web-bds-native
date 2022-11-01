package web

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/bmizerany/pat"
)

func (a *router) adminrender(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := a.be[name]
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

func (a *router) isadmin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := a.session.GetInt(r, "user")
		if id == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := a.App.Users.ID(fmt.Sprint(id))
		if err != nil {
			log.Println(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var isAdmin bool
		for _, role := range user.Roles {
			if role == "admin" {
				isAdmin = true
			}
		}

		if !isAdmin {
			log.Println("user do not have role admin")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *router) adminHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}

func (a *router) adminProducts(w http.ResponseWriter, r *http.Request) {
	p := a.App.AdminProducts.Pagination.Query(r.URL)

	products, err := a.App.Products.Find()
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", 400)
		return
	}
	a.adminrender(w, r, "products.page.html", &templateData{
		Products:   products,
		Pagination: p,
	})
}

func (a *router) adminUsers(w http.ResponseWriter, r *http.Request) {
	p := a.App.AdminUsers.Pagination.Query(r.URL)

	users, err := a.App.AdminUsers.Find()
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", 400)
		return
	}
	a.adminrender(w, r, "users.page.html", &templateData{
		Users:      users,
		Pagination: p,
	})
}

func (a *router) adminUsersDetail(w http.ResponseWriter, r *http.Request) {
	user, err := a.App.AdminUsers.ID(r.URL.Query().Get(":id"))
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", 400)
		return
	}

	a.adminrender(w, r, "users.detail.page.html", &templateData{
		User: user,
		Logs: nil,
	})
}

func (a *router) adminCreateProduct(w http.ResponseWriter, r *http.Request) {
	product, err := a.App.Products.Create()
	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/products/%d/update", product.ID), http.StatusSeeOther)
}

func (a *router) adminUpdateProduct(w http.ResponseWriter, r *http.Request) {
	product, err := a.App.Products.ID(r.URL.Query().Get(":id"))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
		return
	}

	f := product.Form()

	ok := false
	defer func() {
		if !ok {
			a.adminrender(w, r, "products.update.page.html", &templateData{
				Form: f,
			})
		}
	}()

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			f.Errors.Add("err", "err_parse_form")
			return
		}

		f.Values = r.PostForm
		f.Required("Title")

		if !f.Valid() {
			log.Println("form invalid", f.Errors)
			// f.Errors.Add("err", "err_invalid_form")
			return
		}

		if err := a.App.Products.Update(product, f); err != nil {
			log.Println(err)
			f.Errors.Add("err", "could_not_update_product")
			return
		}

		ok = true
		http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
	}

}

func registerAdminRoute(mux *pat.PatternServeMux, a *router) {
	mux.Get("/admin", use(a.adminHome, a.isadmin))
	mux.Get("/admin/products", use(a.adminProducts, a.isadmin))
	mux.Post("/admin/products/:id/update", use(a.adminUpdateProduct, a.isadmin))
	mux.Get("/admin/products/:id/update", use(a.adminUpdateProduct, a.isadmin))
	mux.Get("/admin/products/create", use(a.adminCreateProduct, a.isadmin))

	mux.Get("/admin/users", use(a.adminUsers, a.isadmin))
	mux.Get("/admin/users/:id/detail", use(a.adminUsersDetail, a.isadmin))
}
