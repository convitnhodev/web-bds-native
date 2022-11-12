package web

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/bmizerany/pat"
	"github.com/deeincom/deeincom/pkg/form"
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
			log.Println(err)
			f.Errors.Add("err", "err_could_not_upload")
			return
		}
		defer file.Close()

		fileName, err := a.App.LocalFile.UploadFile(fmt.Sprintf("products/%d/", attachment.Product.ID), file, handler)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		link := a.App.Config.CDNRoot + *fileName

		f.Set("Size", fmt.Sprint(handler.Size))
		f.Set("Link", link)

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

func (a *router) adminUpdateProductMedia(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	productId := query.Get(":id")
	attachmentId := query.Get(":attachmentId")
	typeMedia := query.Get("typeMedia")

	attachment, err := a.App.Attachments.ID(attachmentId)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	key := ""

	switch typeMedia {
	case "poster":
		key = "poster_link"
	case "houseCertificate":
		key = "house_certificate_link"
	case "financePlan":
		key = "finance_plan_link"
	default:
		http.Error(w, "500 - internal server error", 500)
		return
	}

	if err := a.App.Products.Set(productId, key, attachment.Link); err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/products/%s/attachments", productId), http.StatusSeeOther)
}

func (a *router) adminPosts(w http.ResponseWriter, r *http.Request) {
	p := a.App.Posts.Pagination.Query(r.URL)

	posts, err := a.App.Posts.Find()
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", 400)
		return
	}

	a.adminrender(w, r, "posts.page.html", &templateData{
		Posts:      posts,
		Pagination: p,
	})
}

func (a *router) adminCreatePost(w http.ResponseWriter, r *http.Request) {
	f := form.Form{}

	ok := false
	defer func() {
		if !ok {
			a.adminrender(w, r, "posts.update.page.html", &templateData{
				Form: &f,
			})
		}
	}()

	if r.Method == "POST" {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			log.Println(err)
			f.Errors.Add("err", "err_parse_form")
			return
		}

		f.Values = r.PostForm
		f.Required("Title")

		if !f.Valid() {
			log.Println("form invalid", f.Errors)
			return
		}

		id := a.session.GetInt(r, "user")
		post, err := a.App.Posts.Create(id, "blog")

		if err != nil {
			log.Println(err)
			f.Errors.Add("err", "could_not_create_post")
			return
		}

		file, handler, err := r.FormFile("Thumbnail")
		if err != nil && http.ErrMissingFile != err {
			log.Println(err)
			f.Errors.Add("err", "err_could_not_upload")
			return
		}

		if file != nil {
			defer file.Close()

			fileName, err := a.App.LocalFile.UploadFile(fmt.Sprintf("posts/%d/", post.ID), file, handler)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			thumbnail := a.App.Config.CDNRoot + *fileName
			f.Set("Thumbnail", thumbnail)
		} else {
			f.Set("Thumbnail", post.Thumbnail)
		}

		if err := a.App.Posts.Update(post, &f); err != nil {
			log.Println(err)
			f.Errors.Add("err", "could_not_create_post")
			return
		}

		ok = true
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
	}
}

func (a *router) adminUpdatePost(w http.ResponseWriter, r *http.Request) {
	post, err := a.App.Posts.ID(r.URL.Query().Get(":id"))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}

	f := post.Form()

	ok := false
	defer func() {
		if !ok {
			a.adminrender(w, r, "posts.update.page.html", &templateData{
				Form: f,
			})
		}
	}()

	if r.Method == "POST" {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			log.Println(err)
			f.Errors.Add("err", "err_parse_form")
			return
		}

		f.Values = r.PostForm
		f.Required("Title")

		if !f.Valid() {
			log.Println("form invalid", f.Errors)
			return
		}

		file, handler, err := r.FormFile("Thumbnail")
		if err != nil && http.ErrMissingFile != err {
			log.Println(err)
			f.Errors.Add("err", "err_could_not_upload")
			return
		}

		if file != nil {
			defer file.Close()

			fileName, err := a.App.LocalFile.UploadFile(fmt.Sprintf("posts/%d/", post.ID), file, handler)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			thumbnail := a.App.Config.CDNRoot + *fileName
			f.Set("Thumbnail", thumbnail)
		} else {
			f.Set("Thumbnail", post.Thumbnail)
		}

		if err := a.App.Posts.Update(post, f); err != nil {
			log.Println(err)
			f.Errors.Add("err", "could_not_update_post")
			return
		}

		ok = true
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
	}

}

func (a *router) adminRemovePost(w http.ResponseWriter, r *http.Request) {
	a.App.Posts.Remove(r.URL.Query().Get(":id"))
	http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
}

func (a *router) adminChangeCensorshipComment(w http.ResponseWriter, r *http.Request) {
	a.App.Comments.ChangeCensorship(r.URL.Query().Get(":id"))
	w.Write([]byte("Ok"))
}

func registerAdminRoute(mux *pat.PatternServeMux, a *router) {
	mux.Get("/admin", use(a.adminHome, a.isadmin))
	mux.Get("/admin/products", use(a.adminProducts, a.isadmin))

	mux.Get("/admin/products/:id/remove", use(a.adminRemoveProduct, a.isadmin))
	mux.Post("/admin/products/:id/update", use(a.adminUpdateProduct, a.isadmin))
	mux.Get("/admin/products/:id/update", use(a.adminUpdateProduct, a.isadmin))
	mux.Get("/admin/products/create", use(a.adminCreateProduct, a.isadmin))

	mux.Get("/admin/products/:id/attachments", use(a.adminAttachments, a.isadmin))
	mux.Get("/admin/products/:id/attachments/:attachmentId/updateMedia", use(a.adminUpdateProductMedia, a.isadmin))
	mux.Post("/admin/attachments/:id/update", use(a.adminUpdateAttachment, a.isadmin))
	mux.Get("/admin/attachments/:id/update", use(a.adminUpdateAttachment, a.isadmin))
	mux.Get("/admin/attachments/create", use(a.adminCreateAttachment, a.isadmin))

	mux.Get("/admin/posts", use(a.adminPosts, a.isadmin))
	mux.Get("/admin/posts/create", use(a.adminCreatePost, a.isadmin))
	mux.Post("/admin/posts/create", use(a.adminCreatePost, a.isadmin))
	mux.Get("/admin/posts/:id/update", use(a.adminUpdatePost, a.isadmin))
	mux.Post("/admin/posts/:id/update", use(a.adminUpdatePost, a.isadmin))
	mux.Get("/admin/posts/:id/remove", use(a.adminRemovePost, a.isadmin))
	mux.Get("/admin/comments/:id/changeCensorship", use(a.adminChangeCensorshipComment, a.isadmin))

	mux.Get("/admin/users", use(a.adminUsers, a.isadmin))
	mux.Get("/admin/users/:id/detail", use(a.adminUsersDetail, a.isadmin))
}
