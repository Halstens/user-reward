package main

import (
	"github.com/gorilla/mux"
	middleware "github.com/user-reward/internal/middlware"
)

func (app *application) routes() *mux.Router {
	pub := mux.NewRouter()
	pub.HandleFunc("/login", app.Login)

	protect := pub.PathPrefix("").Subrouter()
	protect.Use(middleware.AuthMiddleware(app.jwt))

	protect.HandleFunc("/users/{id}/status", app.ShowUserInfo)
	protect.HandleFunc("/users/leaderboard", app.ShowTopUserByBalance)
	protect.HandleFunc("/users/{id}/task/complete", app.CompletedTask)
	protect.HandleFunc("/users/{id}/referrer", app.AddRefferer)

	return pub
}
