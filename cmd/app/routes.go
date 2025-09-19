package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/{id}/status", app.ShowUserInfo)
	mux.HandleFunc("/users/leaderboard", app.ShowTopUserByBalance)
	mux.HandleFunc("/users/{id}/task/complete", app.CompletedTask)
	mux.HandleFunc("/users/{id}/referrer", app.AddRefferer)

	return mux
}
