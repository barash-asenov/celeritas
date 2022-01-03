package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"myapp/data"
	"net/http"
	"time"
)

func (h *Handlers) UserLogin(rw http.ResponseWriter, r *http.Request) {
	err := h.render(rw, r, "login", nil, nil)

	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) PostUserLogin(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := h.Models.Users.GetByEmail(email)

	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}

	matches, err := user.PasswordMatches(password)

	if err != nil {
		rw.Write([]byte("Error validating password"))
		return
	}

	if !matches {
		rw.Write([]byte("Error validating password"))
		return
	}

	// did the user check rememberMe?
	if r.Form.Get("remember") == "remember" {
		randomString := h.randomString(12)
		hasher := sha256.New()
		_, err := hasher.Write([]byte(randomString))
		if err != nil {
			h.App.ErrorStatus(rw, http.StatusBadRequest)
			return
		}

		sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		rm := data.RememberToken{}
		err = rm.InsertToken(user.ID, sha)
		if err != nil {
			h.App.ErrorStatus(rw, http.StatusBadRequest)
			return
		}

		// set cookie
		expire := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{
			Name: fmt.Sprintf("_%s_remember", h.App.AppName),
			Value: fmt.Sprintf("%d|%s", user.ID, sha),
			Path: "/",
			Expires: expire,
			HttpOnly: true,
			Domain: h.App.Session.Cookie.Domain,
			MaxAge: 315350000,
			Secure: h.App.Session.Cookie.Secure,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(rw, &cookie)

		// save hash in session
		h.App.Session.Put(r.Context(), "remember_token", sha)
	}

	h.App.Session.Put(r.Context(), "userID", user.ID)

	http.Redirect(rw, r, "/", http.StatusSeeOther)
}

func (h *Handlers) Logout(rw http.ResponseWriter, r *http.Request) {
	// delete the remember token if exists
	if h.App.Session.Exists(r.Context(), "remember_token") {
		rt := data.RememberToken{}
		_ = rt.Delete(h.App.Session.GetString(r.Context(), "remember_token"))
	}

	// delete the cookie
	newCookie := http.Cookie{
		Name: fmt.Sprintf("_%s_remember", h.App.AppName),
		Value: "",
		Path: "/",
		Expires: time.Now().Add(-100 * time.Hour),
		HttpOnly: true,
		Domain: h.App.Session.Cookie.Domain,
		MaxAge: -1,
		Secure: h.App.Session.Cookie.Secure,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(rw, &newCookie)

	_ = h.App.Session.RenewToken(r.Context())
	h.App.Session.Remove(r.Context(), "userID")
	h.App.Session.Remove(r.Context(), "rememberToken")
	_ = h.App.Session.Destroy(r.Context())
	_ = h.App.Session.RenewToken(r.Context())

	http.Redirect(rw, r, "/users/login", http.StatusSeeOther)
}
