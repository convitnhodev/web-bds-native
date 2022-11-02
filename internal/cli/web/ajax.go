package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bmizerany/pat"
	"github.com/deeincom/deeincom/pkg/email"
)

func (a *router) ajax(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do not cache ajax
		w.Header().Set("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}

func (a *router) ajaxSendVerifyEmail(w http.ResponseWriter, r *http.Request) {
	// always return ok, no matter what
	defer func() {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	}()

	// user phải logged mới dc yêu cầu verify
	// tránh spam
	id := a.session.GetInt(r, "user")
	if id == 0 {
		log.Println("to request verification email, please logged first")
		return
	}

	user, err := a.App.Users.ID(fmt.Sprint(id))
	if err != nil {
		log.Println(err)
		return
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

func (a *router) ajaxSendVerifyPhone(w http.ResponseWriter, r *http.Request) {
	// always return ok, no matter what
	defer func() {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	}()

	// user phải logged mới dc yêu cầu verify
	// tránh spam
	id := a.session.GetInt(r, "user")
	if id == 0 {
		log.Println("to request verification sms, please logged first")
		return
	}

	user, err := a.App.Users.ID(fmt.Sprint(id))
	if err != nil {
		log.Println(err)
		return
	}
	// nếu user đã verified_phone rồi thì ko cần gởi nữa
	for _, role := range user.Roles {
		if role == "verified_phone" {
			return
		}
	}

	// kiểm tra coi lần trước user yêu cầu verify phone là lúc nào
	// mổi lần yêu cầu phải cách 60s
	if !user.SendVerifiedPhoneAt.Add(60 * time.Second).UTC().Before(time.Now().UTC()) {
		return
	}

	// ghi nhớ lần gởi sms này
	user.Phone = r.URL.Query().Get("phone")
	if err := a.App.Users.LogSendVerifyPhone(user); err != nil {
		log.Println(err)
		return
	}

	// gởi phone verify
}

func registerAjaxRoute(mux *pat.PatternServeMux, a *router) {
	mux.Get("/ajax/verify-email", use(a.ajaxSendVerifyEmail, a.ajax))
	mux.Get("/ajax/verify-phone", use(a.ajaxSendVerifyPhone, a.ajax))
}
