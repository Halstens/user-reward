package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/{id}/status", app.ShowUserInfo)
	mux.HandleFunc("/users/leaderboard", app.)
	mux.HandleFunc("/users/{id}/task/complete ")
	mux.HandleFunc("/users/{id}/referrer")

	return mux
}
