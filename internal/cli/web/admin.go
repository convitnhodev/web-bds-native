package web

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func (a *router) adminRemoveProduct(w http.ResponseWriter, r *http.Request) {
	a.App.AdminProducts.Remove(r.URL.Query().Get(":id"))

	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}

func (a *router) adminHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}

func (a *router) adminAttachments(w http.ResponseWriter, r *http.Request) {
	p := a.App.AdminAttachments.Pagination.Query(r.URL)

	product, err := a.App.Products.ID(r.URL.Query().Get(":id"))
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", 400)
		return
	}
	attachments, err := a.App.AdminAttachments.Find(product)
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", 400)
		return
	}

	a.adminrender(w, r, "attachments.page.html", &templateData{
		Product:     product,
		Pagination:  p,
		Attachments: attachments,
	})
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

func (a *router) adminCreateAttachment(w http.ResponseWriter, r *http.Request) {
	atype := r.URL.Query().Get("type")
	pid := r.URL.Query().Get("pid")
	product, err := a.App.Products.ID(pid)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}
	attachment, err := a.App.Attachments.Create(product, atype)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/attachments/%d/update", attachment.ID), http.StatusSeeOther)
}

func (a *router) adminUpdateAttachment(w http.ResponseWriter, r *http.Request) {
	attachment, err := a.App.Attachments.ID(r.URL.Query().Get(":id"))
	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	product, err := a.App.Products.ID(fmt.Sprint(attachment.Product.ID))
	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	f := attachment.Form()

	ok := false
	defer func() {
		if !ok {
			a.adminrender(w, r, "attachments.update.page.html", &templateData{
				Form:    f,
				Product: product,
			})
		}
	}()

	if r.Method == "POST" {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			f.Errors.Add("err", "err_parse_form")
			return
		}

		f.Values = r.PostForm
		f.Required("Title")

		if !f.Valid() {
			log.Println("form invalid", f.Errors)
			return
		}

		file, handler, err := r.FormFile("UploadFile")
		if err != nil {
			f.Errors.Add("err", "err_could_not_upload")
			return
		}
		defer file.Close()

		f.Set("Size", fmt.Sprint(handler.Size))

		// Create file
		dst, err := os.Create(filepath.Join("upload", handler.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file to the created file on the filesystem
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f.Set("Link", "https://cdn.deein.com/"+handler.Filename)

		if err := a.App.Attachments.Update(attachment, f); err != nil {
			log.Println(err)
			f.Errors.Add("err", "could_not_update_attachment")
			return
		}

		ok = true
		http.Redirect(w, r, fmt.Sprintf("/admin/products/%d/attachments", attachment.Product.ID), http.StatusSeeOther)
	}

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

	mux.Get("/admin/products/:id/remove", use(a.adminRemoveProduct, a.isadmin))
	mux.Post("/admin/products/:id/update", use(a.adminUpdateProduct, a.isadmin))
	mux.Get("/admin/products/:id/update", use(a.adminUpdateProduct, a.isadmin))
	mux.Get("/admin/products/create", use(a.adminCreateProduct, a.isadmin))

	mux.Get("/admin/products/:id/attachments", use(a.adminAttachments, a.isadmin))
	mux.Post("/admin/attachments/:id/update", use(a.adminUpdateAttachment, a.isadmin))
	mux.Get("/admin/attachments/:id/update", use(a.adminUpdateAttachment, a.isadmin))
	mux.Get("/admin/attachments/create", use(a.adminCreateAttachment, a.isadmin))

	mux.Get("/admin/users", use(a.adminUsers, a.isadmin))
	mux.Get("/admin/users/:id/detail", use(a.adminUsersDetail, a.isadmin))
}
