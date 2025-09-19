package main

import (
	"github.com/gorilla/mux"
)

// func (app *application) routes() *http.ServeMux {
// 	// mux := http.NewServeMux()
// 	mux := mux.NewRouter()
// 	mux.HandleFunc("/users/{id}/status", app.ShowUserInfo)
// 	mux.HandleFunc("/users/leaderboard", app.ShowTopUserByBalance)
// 	mux.HandleFunc("/users/{id}/task/complete", app.CompletedTask)
// 	mux.HandleFunc("/users/{id}/referrer", app.AddRefferer)

// 	return mux
// }

func (app *application) routes() *mux.Router {
	m := mux.NewRouter()
	m.HandleFunc("/users/{id}/status", app.ShowUserInfo)
	m.HandleFunc("/users/leaderboard", app.ShowTopUserByBalance)
	m.HandleFunc("/users/{id}/task/complete", app.CompletedTask)
	m.HandleFunc("/users/{id}/referrer", app.AddRefferer)

	return m
}
