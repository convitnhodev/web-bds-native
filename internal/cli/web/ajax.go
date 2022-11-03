package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bmizerany/pat"
	"github.com/deeincom/deeincom/pkg/email"
	"github.com/deeincom/deeincom/pkg/models"
	"github.com/deeincom/deeincom/pkg/phone"
)

func (a *router) ajax(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do not cache ajax
		w.Header().Set("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}

// gởi verify email cho user
func (a *router) ajaxSendVerifyEmail(w http.ResponseWriter, r *http.Request) {
	// luôn ok
	defer func() {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	}()

	emailladdress := r.URL.Query().Get("email")
	// lấy user
	var user *models.User
	// nếu user ko đăng nhập, lấy user bằng phone
	id := a.session.GetInt(r, "user")
	if id == 0 {
		if u, err := a.App.Users.GetByEmail(emailladdress); err != nil {
			log.Println(err)
			return
		} else {
			user = u
		}
	} else {
		if u, err := a.App.Users.ID(fmt.Sprint(id)); err != nil {
			log.Println(err)
			return
		} else {
			user = u
		}
	}

	// nếu user đã verified_email rồi thì ko cần gởi nữa
	for _, role := range user.Roles {
		if role == "verified_email" {
			return
		}
	}

	// kiểm tra coi lần trước user yêu cầu verify email là lúc nào
	// mổi lần yêu cầu phải cách 60s
	if !user.SendVerifiedEmailAt.Add(60 * time.Second).UTC().Before(time.Now().UTC()) {
		return
	}

	// ghi nhớ lần gởi email này
	user.Email = r.URL.Query().Get("email")
	if err := a.App.Users.LogSendVerifyEmail(user); err != nil {
		log.Println(err)
		return
	}

	// gởi email verify
	email.PostmarkApiToken = a.App.Config.PostmarkApiToken
	if err := email.SendVerifyEmail(user); err != nil {
		log.Println(err)
		return
	}

}

// gởi một tin nhắn sms để verify phone
func (a *router) ajaxSendVerifyPhone(w http.ResponseWriter, r *http.Request) {
	// luôn ok
	defer func() {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	}()

	phonenumber := r.URL.Query().Get("phone")
	// lấy user
	var user *models.User
	// nếu user ko đăng nhập, lấy user bằng phone
	id := a.session.GetInt(r, "user")
	if id == 0 {
		if u, err := a.App.Users.GetByPhone(phonenumber); err != nil {
			log.Println(err)
			return
		} else {
			user = u
		}
	} else {
		if u, err := a.App.Users.ID(fmt.Sprint(id)); err != nil {
			log.Println(err)
			return
		} else {
			user = u
		}
	}

	// nếu user đã verified_phone rồi thì ko cần gởi nữa
	for _, role := range user.Roles {
		if role == "verified_phone" {
			return
		}
	}

	// mổi lần yêu cầu phải cách 60s
	if !user.SendVerifiedPhoneAt.Add(60 * time.Second).UTC().Before(time.Now().UTC()) {
		return
	}

	// ghi nhớ lần gởi sms này
	if err := a.App.Users.LogSendVerifyPhone(user); err != nil {
		log.Println(err)
		return
	}

	// gởi phone verify
	phone.ESMS_APIKEY = a.App.Config.ESMS_APIKEY
	phone.ESMS_SECRET = a.App.Config.ESMS_SECRET
	if err := phone.SendSMS(user); err != nil {
		log.Println(err)
	}
}

func registerAjaxRoute(mux *pat.PatternServeMux, a *router) {
	mux.Get("/ajax/verify-email", use(a.ajaxSendVerifyEmail, a.ajax))
	mux.Get("/ajax/verify-phone", use(a.ajaxSendVerifyPhone, a.ajax))
}
