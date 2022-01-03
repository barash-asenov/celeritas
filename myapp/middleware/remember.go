package middleware

import (
	"fmt"
	"myapp/data"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (m *Middleware) CheckRemember(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if !m.App.Session.Exists(r.Context(), "userID") {
			// user is not logged in
			cookie, err := r.Cookie(fmt.Sprintf("_%s_remember", m.App.AppName))
			if err != nil {
				// no cookie, so on to the next middleware
				next.ServeHTTP(rw, r)
			} else {
				// we found a cookie, so check it
				key := cookie.Value
				var u data.User
				if len(key) > 0 {
					// cookie has some data, validate it
					split := strings.Split(key, "|")
					uid, hash := split[0], split[1]
					id, _ := strconv.Atoi(uid)
					validHash := u.CheckForRememberToken(id, hash)
					if !validHash {
						m.deleteRememberCookie(rw, r)
						m.App.Session.Put(r.Context(), "error", "You've been logged out from another device")
						next.ServeHTTP(rw, r)
					} else {
						// valid hash, so log the user in
						user, _ := u.Get(id)
						m.App.Session.Put(r.Context(), "userID", user.ID)
						m.App.Session.Put(r.Context(), "rememberToken", hash)
						next.ServeHTTP(rw, r)
					}
				} else {
					// key length is zero, probably leftover cookie (user has not closed browser)
					m.deleteRememberCookie(rw, r)
					next.ServeHTTP(rw, r)
				}
			}
		} else {
			// user is logged in
			next.ServeHTTP(rw, r)
		}
	})
}

func (m *Middleware) deleteRememberCookie(rw http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	// delete the cookie
	newCookie := http.Cookie{
		Name: fmt.Sprintf("_%s_remember", m.App.AppName),
		Value: "",
		Path: "/",
		Expires: time.Now().Add(-100 * time.Hour),
		HttpOnly: true,
		Domain: m.App.Session.Cookie.Domain,
		MaxAge: -1,
		Secure: m.App.Session.Cookie.Secure,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(rw, &newCookie)

	// log the user out
	m.App.Session.Remove(r.Context(), "userID")
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
}