package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/bmizerany/pat"
	"github.com/deeincom/deeincom/pkg/appotapay"
	"github.com/deeincom/deeincom/pkg/files"
	"github.com/deeincom/deeincom/pkg/form"
	"github.com/deeincom/deeincom/pkg/models"
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

func (a *router) ispartner(next http.HandlerFunc) http.HandlerFunc {
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

		var isPartner bool
		var isAdmin bool

		for _, role := range user.Roles {
			if role == "admin" {
				isAdmin = true
			}

			if role == "deein_partner" {
				isPartner = true
			}
		}

		if !isAdmin && !isPartner {
			log.Println("user do not have role admin")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		role := "deein_partner"
		if isAdmin {
			role = "admin"
		}

		ctx := context.WithValue(r.Context(), "role", role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *router) adminRemoveProduct(w http.ResponseWriter, r *http.Request) {
	userId := a.session.Get(r, "user")
	productId := r.URL.Query().Get(":id")
	a.App.AdminProducts.Remove(productId)

	a.App.Log.Add(fmt.Sprint(userId), fmt.Sprintf("Người dùng %d xoá sản phẩm %s.", userId, productId))
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
	userId := a.session.GetInt(r, "user")
	role := r.Context().Value("role")

	var products []*models.Product
	var err error
	if role == "admin" {
		products, err = a.App.Products.Find()
	} else {
		products, err = a.App.Products.CreatedBy(userId)
	}

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
	kycStatus := r.URL.Query().Get("kyc_status")
	partnerStatus := r.URL.Query().Get("partner_status")

	users, err := a.App.AdminUsers.Find(kycStatus, partnerStatus)
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", 400)
		return
	}

	a.adminrender(w, r, "users.page.html", &templateData{
		Users:          users,
		IsKYCQuery:     kycStatus != "",
		IsPartnerQuery: partnerStatus != "",
		Pagination:     p,
	})
}

func (a *router) adminUsersDetail(w http.ResponseWriter, r *http.Request) {
	user, err := a.App.AdminUsers.ID(r.URL.Query().Get(":id"))

	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", 400)
		return
	}

	kycList, err := a.App.KYC.User(r.URL.Query().Get(":id"))
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", 400)
		return
	}

	partnerList, err := a.App.Partner.User(r.URL.Query().Get(":id"))
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", 400)
		return
	}

	a.adminrender(w, r, "users.detail.page.html", &templateData{
		UserInfo:    user,
		KYCList:     kycList,
		PartnerList: partnerList,
		Logs:        nil,
	})
}

func (a *router) adminApproveKYC(w http.ResponseWriter, r *http.Request) {
	approvedBy := a.session.GetInt(r, "user")
	userId := r.URL.Query().Get(":id")
	kycId := r.URL.Query().Get(":kycId")

	// Update KYC status
	status := "approved_kyc"
	err := a.App.KYC.FeedbackKYC(
		kycId,
		fmt.Sprint(approvedBy),
		status,
		"",
	)

	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	// Update Role, KCS status user
	err = a.App.Users.UpdateKYCStatus(
		userId,
		status,
	)

	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	a.App.Log.Add(fmt.Sprint(approvedBy), fmt.Sprintf("Người dùng %d đã đồng ý KYC người dùng %s ở yêu cầu %s.", approvedBy, userId, kycId))
	http.Redirect(w, r, fmt.Sprintf("/admin/users/%s/detail", userId), http.StatusSeeOther)
}

func (a *router) adminRejectKYC(w http.ResponseWriter, r *http.Request) {
	rejectedBy := a.session.GetInt(r, "user")
	userId := r.URL.Query().Get(":id")
	kycId := r.URL.Query().Get(":kycId")
	ok := false
	f := form.New(nil)

	defer func() {
		if !ok {
			user, err := a.App.AdminUsers.ID(r.URL.Query().Get(":id"))
			if err != nil {
				log.Println(err)
				return
			}

			kycList, err := a.App.KYC.User(r.URL.Query().Get(":id"))
			if err != nil {
				log.Println(err)
				return
			}

			partnerList, err := a.App.Partner.User(r.URL.Query().Get(":id"))
			if err != nil {
				log.Println(err)
				return
			}

			a.adminrender(w, r, "users.detail.page.html", &templateData{
				Form:        f,
				UserInfo:    user,
				PartnerList: partnerList,
				KYCList:     kycList,
			})
		}
	}()

	if err := r.ParseForm(); err != nil {
		f.Errors.Add("err", "err_parse_form")
		return
	}

	f.Values = r.PostForm

	// Update KYC status
	status := "rejected_kyc"
	err := a.App.KYC.FeedbackKYC(
		kycId,
		fmt.Sprint(rejectedBy),
		status,
		f.Get("Feedback"),
	)

	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	// Update Role, KCS status user
	err = a.App.Users.UpdateKYCStatus(
		userId,
		status,
	)

	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	ok = true
	a.App.Log.Add(fmt.Sprint(rejectedBy), fmt.Sprintf("Người dùng %d đã từ chối KYC người dùng %s ở yêu cầu %s.", rejectedBy, userId, kycId))
	http.Redirect(w, r, fmt.Sprintf("/admin/users/%s/detail", userId), http.StatusSeeOther)
}

func (a *router) adminApprovePartner(w http.ResponseWriter, r *http.Request) {
	approvedBy := a.session.GetInt(r, "user")
	userId := r.URL.Query().Get(":id")
	partnerId := r.URL.Query().Get(":partnerId")

	// Update KYC status
	status := "approved"
	err := a.App.Partner.FeedbackPartner(
		partnerId,
		fmt.Sprint(approvedBy),
		status,
		"",
	)

	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	// Update Role, KCS status user
	err = a.App.Users.UpdatePartnerStatus(
		userId,
		status,
	)

	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	a.App.Log.Add(fmt.Sprint(approvedBy), fmt.Sprintf("Người dùng %d đã đồng ý người dùng %s thành đối tác ở yêu cầu %s.", approvedBy, userId, partnerId))
	http.Redirect(w, r, fmt.Sprintf("/admin/users/%s/detail", userId), http.StatusSeeOther)
}

func (a *router) adminRejectPartner(w http.ResponseWriter, r *http.Request) {
	rejectedBy := a.session.GetInt(r, "user")
	userId := r.URL.Query().Get(":id")
	partnerId := r.URL.Query().Get(":partnerId")
	ok := false
	f := form.New(nil)

	defer func() {
		if !ok {
			user, err := a.App.AdminUsers.ID(r.URL.Query().Get(":id"))
			if err != nil {
				log.Println(err)
				return
			}

			kycList, err := a.App.KYC.User(r.URL.Query().Get(":id"))
			if err != nil {
				log.Println(err)
				return
			}

			partnerList, err := a.App.Partner.User(r.URL.Query().Get(":id"))
			if err != nil {
				log.Println(err)
				return
			}

			a.adminrender(w, r, "users.detail.page.html", &templateData{
				Form:        f,
				UserInfo:    user,
				PartnerList: partnerList,
				KYCList:     kycList,
			})
		}
	}()

	if err := r.ParseForm(); err != nil {
		f.Errors.Add("err", "err_parse_form")
		return
	}

	f.Values = r.PostForm

	// Update KYC status
	status := "rejected"
	err := a.App.Partner.FeedbackPartner(
		partnerId,
		fmt.Sprint(rejectedBy),
		status,
		f.Get("Feedback"),
	)

	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	// Update Role, KCS status user
	err = a.App.Users.UpdatePartnerStatus(
		userId,
		status,
	)

	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	ok = true

	a.App.Log.Add(fmt.Sprint(rejectedBy), fmt.Sprintf("Người dùng %d đã từ chối người dùng %s thành đối tác ở yêu cầu %s.", rejectedBy, userId, partnerId))
	http.Redirect(w, r, fmt.Sprintf("/admin/users/%s/detail", userId), http.StatusSeeOther)
}

func (a *router) adminCreateProduct(w http.ResponseWriter, r *http.Request) {
	userId := a.session.GetInt(r, "user")
	role := r.Context().Value("role")

	var isCensorship bool
	if role == "admin" {
		isCensorship = true
	} else {
		isCensorship = false
	}

	product, err := a.App.Products.Create(userId, isCensorship)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	a.App.Log.Add(fmt.Sprint(userId), fmt.Sprintf("Người dùng %d tạo sản phẩm %d.", userId, product.ID))
	http.Redirect(w, r, fmt.Sprintf("/admin/products/%d/update", product.ID), http.StatusSeeOther)
}

func (a *router) adminApproveProduct(w http.ResponseWriter, r *http.Request) {
	userId := a.session.GetInt(r, "user")
	productId := r.URL.Query().Get(":id")

	err := a.App.Products.Approve(productId)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	a.App.Log.Add(fmt.Sprint(userId), fmt.Sprintf("Người dùng %d duyệt sản phẩm %s.", userId, productId))
	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}

func (a *router) adminCreateAttachment(w http.ResponseWriter, r *http.Request) {
	userId := a.session.GetInt(r, "user")
	atype := r.URL.Query().Get("type")
	pid := r.URL.Query().Get(":id")

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

	a.App.Log.Add(fmt.Sprint(userId), fmt.Sprintf("Người dùng %d tạo tệp đính kèm %d cho sản phẩm %d", userId, attachment.ID, product.ID))
	http.Redirect(w, r, fmt.Sprintf("/admin/products/%s/attachments/%d/update", pid, attachment.ID), http.StatusSeeOther)
}

func (a *router) adminUpdateAttachment(w http.ResponseWriter, r *http.Request) {
	userId := a.session.GetInt(r, "user")
	attachment, err := a.App.Attachments.ID(r.URL.Query().Get(":attachmentId"))
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
		r.Body = http.MaxBytesReader(w, r.Body, 101<<20)
		if err := r.ParseMultipartForm(101 << 20); err != nil {
			f.Errors.Add("err", "err_parse_form")
			return
		}

		f.Values = r.PostForm
		f.Required("Title")

		if !f.Valid() {
			log.Println("form invalid", f.Errors)
			return
		}

		uploadFileLink := ""
		if attachment.Link == "" {
			file, handler, err := r.FormFile("UploadFile")

			if err != nil {
				log.Println(err)
				f.Errors.Add("err", "err_could_not_upload")
				return
			}
			defer file.Close()

			fileName, err := a.App.LocalFile.UploadFile(fmt.Sprintf("products.%d/", attachment.Product.ID), file, handler)
			if err != nil {
				log.Println(err)
				if errors.Is(err, files.FileExists) {
					f.Errors.Add("err", "could_not_attachment_exists")
				} else {
					f.Errors.Add("err", "err_could_not_upload")
				}
				return
			}

			f.Set("Size", fmt.Sprint(handler.Size))
			f.Set("Link", *fileName)
			uploadFileLink = *fileName
		}

		if err := a.App.Attachments.Update(attachment, f); err != nil {
			log.Println(err)
			f.Errors.Add("err", "could_not_update_attachment")
			// Remove file khi không update được Attachments
			if uploadFileLink != "" {
				a.App.LocalFile.RemoveLocalFile(uploadFileLink)
			}
			return
		}

		ok = true

		a.App.Log.Add(fmt.Sprint(userId), fmt.Sprintf("Người dùng %d cập nhật thông tin tệp đính kèm %d cho sản phẩm %d.", userId, attachment.ID, product.ID))
		http.Redirect(w, r, fmt.Sprintf("/admin/products/%d/attachments", attachment.Product.ID), http.StatusSeeOther)
	}
}

func (a *router) adminRemoveAttchment(w http.ResponseWriter, r *http.Request) {
	userId := a.session.GetInt(r, "user")
	productId := r.URL.Query().Get(":id")
	attachmentId := r.URL.Query().Get(":attachmentId")

	product, err := a.App.Products.ID(productId)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
		return
	}

	attachment, err := a.App.Attachments.ID(attachmentId)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	// unlink docs
	if product.FinancePlanLink == attachment.Link {
		if err := a.App.Products.Set(productId, "finance_plan_link", ""); err != nil {
			log.Println(err)
			http.Error(w, "500 - internal server error", 500)
			return
		}
	}

	// unlink docs
	if product.HouseCertificateLink == attachment.Link {
		if err := a.App.Products.Set(productId, "house_certificate_link", ""); err != nil {
			log.Println(err)
			http.Error(w, "500 - internal server error", 500)
			return
		}
	}

	// unlink docs
	if product.PosterLink == attachment.Link {
		if err := a.App.Products.Set(productId, "poster_link", ""); err != nil {
			log.Println(err)
			http.Error(w, "500 - internal server error", 500)
			return
		}
	}

	err = a.App.Attachments.Remove(attachmentId)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	a.App.LocalFile.RemoveLocalFile(attachment.Link)
	a.App.Log.Add(fmt.Sprint(userId), fmt.Sprintf("Người dùng %d xoá tệp đính kèm %s ở sản phẩm %s.", userId, attachmentId, productId))

	http.Redirect(w, r, fmt.Sprintf("/admin/products/%s/attachments", productId), http.StatusSeeOther)
}

func (a *router) adminUpdateProduct(w http.ResponseWriter, r *http.Request) {
	userId := a.session.GetInt(r, "user")
	product, err := a.App.Products.ID(r.URL.Query().Get(":id"))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
		return
	}

	countInvoice, err := a.App.InvoiceItem.CountByProduct(product.ID)
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

		// Khi có hoá đơn không cho sửa slot
		if *countInvoice > 0 && f.GetInt("NumOfSlot") != product.NumOfSlot {
			f.Errors.Add("NumOfSlot", "err_product_buying")
			return
		}

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

		a.App.Log.Add(fmt.Sprint(userId), fmt.Sprintf("Người dùng %d cập nhật thông tin sản phẩm %d.", userId, product.ID))
		http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
	}
}

func (a *router) adminEnableSellingProduct(w http.ResponseWriter, r *http.Request) {
	userId := a.session.GetInt(r, "user")
	productId := r.URL.Query().Get(":id")

	product, err := a.App.Products.ID(productId)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
		return
	}

	err = a.App.Products.EnableSelling(productId)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
		return
	}

	log := fmt.Sprintf("Người dùng %d mở mua bán cho sản phẩm %d.", userId, product.ID)
	if product.IsSelling {
		log = fmt.Sprintf("Người dùng %d tắt mua bán cho sản phẩm %d.", userId, product.ID)
	}

	a.App.Log.Add(fmt.Sprint(userId), log)
	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}

func (a *router) adminUpdateProductMedia(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	productId := query.Get(":id")
	attachmentId := query.Get(":attachmentId")
	typeMedia := query.Get("typeMedia")
	userId := a.session.GetInt(r, "user")

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

	a.App.Log.Add(fmt.Sprint(userId), fmt.Sprintf("Người dùng %d cập nhật %s cho sản phẩm %s.", userId, typeMedia, productId))

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
	f := form.New(nil)

	ok := false
	defer func() {
		if !ok {
			a.adminrender(w, r, "posts.update.page.html", &templateData{
				Form: f,
			})
		}
	}()

	if r.Method == "POST" {
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
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

			fileName, err := a.App.LocalFile.UploadFile(fmt.Sprintf("posts.%d/", post.ID), file, handler)
			if err != nil && !errors.Is(err, files.FileExists) {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			f.Set("Thumbnail", *fileName)
		} else {
			f.Set("Thumbnail", post.Thumbnail)
		}

		if err := a.App.Posts.Update(post, f); err != nil {
			log.Println(err)
			f.Errors.Add("err", "could_not_create_post")
			return
		}

		ok = true
		a.App.Log.Add(fmt.Sprint(id), fmt.Sprintf("Người dùng %d đăng bài viết mới %d.", id, post.ID))

		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
	}
}

func (a *router) adminUpdatePost(w http.ResponseWriter, r *http.Request) {
	userId := a.session.GetInt(r, "user")
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
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
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

			fileName, err := a.App.LocalFile.UploadFile(fmt.Sprintf("posts.%d/", post.ID), file, handler)
			if err != nil && !errors.Is(err, files.FileExists) {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			f.Set("Thumbnail", *fileName)
		} else {
			f.Set("Thumbnail", post.Thumbnail)
		}

		if err := a.App.Posts.Update(post, f); err != nil {
			log.Println(err)
			f.Errors.Add("err", "could_not_update_post")
			return
		}

		ok = true
		a.App.Log.Add(fmt.Sprint(userId), fmt.Sprintf("Người dùng %d cập nhật bài viết %d.", userId, post.ID))

		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
	}
}

func (a *router) adminRemovePost(w http.ResponseWriter, r *http.Request) {
	userId := a.session.GetInt(r, "user")
	postId := r.URL.Query().Get(":id")

	a.App.Posts.Remove(postId)
	a.App.Log.Add(fmt.Sprint(userId), fmt.Sprintf("Người dùng %d xoá bài viết %s.", userId, postId))

	http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
}

func (a *router) adminChangeCensorshipComment(w http.ResponseWriter, r *http.Request) {
	userId := a.session.GetInt(r, "user")
	commentId := r.URL.Query().Get(":id")

	a.App.Comments.ChangeCensorship(commentId)
	a.App.Log.Add(fmt.Sprint(userId), fmt.Sprintf("Người dùng %d duyệt bình luận %s.", userId, commentId))

	http.Redirect(w, r, "/admin/comments", http.StatusSeeOther)
}

func (a *router) adminComments(w http.ResponseWriter, r *http.Request) {
	p := a.App.Comments.Pagination.Query(r.URL)

	comments, err := a.App.Comments.Find()
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", 400)
		return
	}

	a.adminrender(w, r, "comments.page.html", &templateData{
		Comments:   comments,
		Pagination: p,
	})
}

func (a *router) adminRemoveComent(w http.ResponseWriter, r *http.Request) {
	userId := a.session.GetInt(r, "user")
	commentId := r.URL.Query().Get(":id")

	a.App.Comments.Remove(commentId)
	a.App.Log.Add(fmt.Sprint(userId), fmt.Sprintf("Người dùng %d xoá bình luận %s.", userId, commentId))

	http.Redirect(w, r, "/admin/comments", http.StatusSeeOther)
}

func (a *router) adminLogs(w http.ResponseWriter, r *http.Request) {
	userInfo := r.URL.Query().Get("user_info")
	date := r.URL.Query().Get("date")
	p := a.App.Log.Pagination.Query(r.URL)
	f := form.New(nil)

	f.Set("Date", date)
	f.Set("UserInfo", userInfo)

	var user *models.User
	var err error
	var logs []*models.Log
	var lerr error

	defer func() {
		a.adminrender(w, r, "logs.page.html", &templateData{
			Logs:       logs,
			Form:       f,
			Pagination: p,
		})
	}()

	if userInfo == "" {
		logs, lerr = a.App.Log.Find("", date)
	} else {
		if strings.Contains(userInfo, "@") {
			user, err = a.App.Users.GetByEmail(userInfo)
		} else {
			user, err = a.App.Users.GetByPhone(userInfo)
		}

		if err != nil {
			log.Println(err)
			f.Errors.Add("err", "user_query_err")
			return
		}

		logs, lerr = a.App.Log.Find(fmt.Sprint(user.ID), date)
	}

	if lerr != nil {
		log.Println(lerr)
		http.Error(w, "bad request", 400)
		return
	}
}

func (a *router) adminInvoices(w http.ResponseWriter, r *http.Request) {
	p := a.App.Invoice.Pagination.Query(r.URL)

	product, err := a.App.Products.ID(r.URL.Query().Get(":id"))
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", 400)
		return
	}

	invoices, err := a.App.Invoice.Find(product.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", 400)
		return
	}

	a.adminrender(w, r, "invoices.page.html", &templateData{
		Product:    product,
		Pagination: p,
		Invoices:   invoices,
	})
}

func (a *router) adminViewInvoice(w http.ResponseWriter, r *http.Request) {
	productId := r.URL.Query().Get(":id")
	invoice, err := a.App.Invoice.ID(r.URL.Query().Get(":invoiceId"))
	if err != nil {
		log.Println(err)
		http.Error(w, "500 - internal server error", 500)
		return
	}

	payments, err := a.App.Payment.InvoiceID(invoice.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", 400)
		return
	}

	invoiceItems, err := a.App.InvoiceItem.InvoiceID(invoice.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", 400)
		return
	}

	sumAmount := 0
	for _, it := range invoiceItems {
		sumAmount += it.Amount
	}

	a.adminrender(w, r, "invoices.view.page.html", &templateData{
		SumAmount:    sumAmount,
		Invoice:      invoice,
		Payments:     payments,
		InvoiceItems: invoiceItems,
		ProductId:    productId,
	})
}

func (a *router) adminRefundInvoice(w http.ResponseWriter, r *http.Request) {
	ProductId := r.URL.Query().Get(":id")
	invoiceId := r.URL.Query().Get(":invoiceId")
	payemntId := r.URL.Query().Get(":paymentId")

	invoice, err := a.App.Invoice.ID(invoiceId)
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", 400)
		return
	}

	payment, err := a.App.Payment.ID(payemntId)
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", 400)
		return
	}

	// Kiểm tra status trước khi refund
	if invoice.Status == "deposit" && payment.Status == "success" {
		appotapay.APTPaymentHost = a.App.Config.APTPaymentHost
		appotapay.ApiKey = a.App.Config.APTApiKey
		appotapay.PartnerCode = a.App.Config.APTPartnerCode
		appotapay.SecretKey = a.App.Config.APTSecretKey

		refundPostData := appotapay.APTRefundPayload{
			Amount:           payment.Amount,
			RefundId:         fmt.Sprintf("refund-%d", payment.ID),
			AppotapayTransId: payment.AppotapayTransId,
			Reason:           fmt.Sprintf("Hoàn tiền cho hoá đơn %d", invoice.ID),
		}
		refundRes, err := appotapay.Refund(refundPostData)

		if err != nil {
			log.Println(err)
			http.Error(w, "bad request", 400)
			return
		}

		refundDataStr, err := json.Marshal(refundRes)

		if err != nil {
			log.Println(err)
			http.Error(w, "bad request", 400)
			return
		}

		// Update thông tin refund và status
		err = a.App.Payment.Refund(
			payment.ID,
			string(refundDataStr),
			refundRes.Data.RefundId,
			refundRes.Data.TransactionTs,
		)
		if err != nil {
			log.Println(err)
			http.Error(w, "bad request", 400)
			return
		}

		err = a.App.Invoice.Refund(invoice.ID)
		if err != nil {
			log.Println(err)
			http.Error(w, "bad request", 400)
			return
		}

		// Update Product: số slot còn lại tăng lên
		invoiceItems, err := a.App.InvoiceItem.InvoiceID(invoice.ID)
		if err != nil {
			log.Println(err)
			http.Error(w, "bad request", 400)
			return
		}

		for _, it := range invoiceItems {
			err = a.App.Products.AddRemain(it.ProductId, it.Quatity)
			if err != nil {
				log.Println(err)
				http.Error(w, "bad request", 400)
				return
			}
		}
	}

	http.Redirect(
		w,
		r,
		fmt.Sprintf("/admin/products/%s/invoices/%s/view", ProductId, invoiceId),
		http.StatusSeeOther,
	)
}

func (a *router) adminCollectMoneyInvoice(w http.ResponseWriter, r *http.Request) {
	// TODO: refund
}

func (a *router) adminCloseInvoice(w http.ResponseWriter, r *http.Request) {
	// TODO: refund
}

func registerAdminRoute(mux *pat.PatternServeMux, a *router) {
	mux.Get("/admin", use(a.adminHome, a.ispartner))
	mux.Get("/admin/products", use(a.adminProducts, a.ispartner))

	mux.Get("/admin/products/:id/remove", use(a.adminRemoveProduct, a.isadmin))
	mux.Post("/admin/products/:id/update", use(a.adminUpdateProduct, a.ispartner))
	mux.Get("/admin/products/:id/update", use(a.adminUpdateProduct, a.ispartner))
	mux.Get("/admin/products/create", use(a.adminCreateProduct, a.ispartner))
	mux.Get("/admin/products/:id/enableSelling", use(a.adminEnableSellingProduct, a.ispartner))

	mux.Get("/admin/products/:id/attachments", use(a.adminAttachments, a.ispartner))
	mux.Get("/admin/products/:id/attachments/:attachmentId/updateMedia", use(a.adminUpdateProductMedia, a.ispartner))
	mux.Post("/admin/products/:id/attachments/:attachmentId/update", use(a.adminUpdateAttachment, a.ispartner))
	mux.Get("/admin/products/:id/attachments/:attachmentId/update", use(a.adminUpdateAttachment, a.ispartner))
	mux.Get("/admin/products/:id/attachments/create", use(a.adminCreateAttachment, a.ispartner))
	mux.Get("/admin/products/:id/attachments/:attachmentId/remove", use(a.adminRemoveAttchment, a.ispartner))
	mux.Get("/admin/products/:id/approve", use(a.adminApproveProduct, a.isadmin))
	mux.Get("/admin/products/:id/invoices", use(a.adminInvoices, a.isadmin))
	mux.Get("/admin/products/:id/invoices/:invoiceId/view", use(a.adminViewInvoice, a.isadmin))

	mux.Get("/admin/products/:id/invoices/:invoiceId/refund/:paymentId", use(a.adminRefundInvoice, a.isadmin))
	mux.Get("/admin/products/:id/invoices/:invoiceId/collectMoney", use(a.adminCollectMoneyInvoice, a.isadmin))
	mux.Get("/admin/products/:id/invoices/:invoiceId/closeInvoice", use(a.adminCloseInvoice, a.isadmin))

	mux.Get("/admin/posts", use(a.adminPosts, a.isadmin))
	mux.Get("/admin/posts/create", use(a.adminCreatePost, a.isadmin))
	mux.Post("/admin/posts/create", use(a.adminCreatePost, a.isadmin))
	mux.Get("/admin/posts/:id/update", use(a.adminUpdatePost, a.isadmin))
	mux.Post("/admin/posts/:id/update", use(a.adminUpdatePost, a.isadmin))
	mux.Get("/admin/posts/:id/remove", use(a.adminRemovePost, a.isadmin))

	mux.Get("/admin/comments", use(a.adminComments, a.isadmin))
	mux.Get("/admin/comments/:id/changeCensorship", use(a.adminChangeCensorshipComment, a.isadmin))
	mux.Get("/admin/comments/:id/remove", use(a.adminRemoveComent, a.isadmin))

	mux.Get("/admin/users", use(a.adminUsers, a.isadmin))
	mux.Get("/admin/users/:id/detail", use(a.adminUsersDetail, a.isadmin))

	mux.Get("/admin/users/:id/kyc/:kycId/approve", use(a.adminApproveKYC, a.isadmin))
	mux.Post("/admin/users/:id/kyc/:kycId/reject", use(a.adminRejectKYC, a.isadmin))

	mux.Get("/admin/users/:id/partner/:partnerId/approve", use(a.adminApprovePartner, a.isadmin))
	mux.Post("/admin/users/:id/partner/:partnerId/reject", use(a.adminRejectPartner, a.isadmin))

	mux.Get("/admin/logs", use(a.adminLogs, a.isadmin))
}
