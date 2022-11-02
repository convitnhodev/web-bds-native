package web

import (
	"encoding/json"
	"net/http"

	"github.com/bmizerany/pat"
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
}

func (a *router) ajaxSendVerifyPhone(w http.ResponseWriter, r *http.Request) {
	// always return ok, no matter what
	defer func() {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	}()
}

func registerAjaxRoute(mux *pat.PatternServeMux, a *router) {
	mux.Get("/ajax/verify-email", use(a.ajaxSendVerifyEmail, a.ajax))
	mux.Get("/ajax/verify-phone", use(a.ajaxSendVerifyPhone, a.ajax))
}
