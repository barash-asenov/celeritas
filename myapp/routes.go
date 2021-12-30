package main

import (
	"fmt"
	"myapp/data"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes

	// add routes here
	a.get("/", a.Handlers.Home)
	a.get("/go-page", a.Handlers.GoPage)
	a.get("/jet-page", a.Handlers.JetPage)
	a.get("/sessions", a.Handlers.SessionTest)

	a.get("/users/login", a.Handlers.UserLogin)
	a.post("/users/login", a.Handlers.PostUserLogin)
	a.get("/users/logout", a.Handlers.Logout)

	a.get("/form", a.Handlers.Form)
	a.post("/form", a.Handlers.PostForm)

	a.get("/json", a.Handlers.JSON)
	a.get("/xml", a.Handlers.XML)
	a.get("/download-file", a.Handlers.DownloadFile)

	a.get("/cache-test", a.Handlers.ShowCachePage)
	a.post("/api/save-in-cache", a.Handlers.SaveInCache)
	a.post("/api/get-from-cache", a.Handlers.GetFromCache)
	a.post("/api/delete-from-cache", a.Handlers.DeleteFromCache)
	a.post("/api/empty-cache", a.Handlers.EmptyCache)

	a.get("/create-user", func(rw http.ResponseWriter, r *http.Request) {
		u := data.User{
			FirstName: "Barash",
			LastName:  "Asenov",
			Email:     "barash.asenov@gmail.com",
			Active:    1,
			Password:  "password",
		}

		id, err := a.Models.Users.Insert(u)

		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		fmt.Fprintf(rw, "%d: %s", id, u.FirstName)
	})

	a.get("/get-all-users", func(rw http.ResponseWriter, r *http.Request) {
		users, err := a.Models.Users.GetAll()

		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		for _, x := range users {
			fmt.Fprintf(rw, x.LastName)
		}
	})

	a.get("/get-user/{id}", func(rw http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		u, err := a.Models.Users.Get(id)

		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		fmt.Fprintf(rw, "%s %s %s", u.FirstName, u.LastName, u.Email)
	})

	a.get("/update-user/{id}", func(rw http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		u, err := a.Models.Users.Get(id)

		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		u.LastName = a.App.RandomString(10)

		validator := a.App.Validator(nil)
		u.LastName = ""

		u.Validate(validator)

		if !validator.Valid() {
			fmt.Fprintf(rw, "failed validation")
			return
		}

		err = u.Update(*u)

		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		fmt.Fprintf(rw, "updated last name to %s", u.LastName)
	})

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.App.Routes
}
