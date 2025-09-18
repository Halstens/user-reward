package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/user-reward/internal/models"
)

func (app *application) ShowUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow: ", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

}

func (app *application) Transaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow: ", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	var request models.WalletOperation
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if request.Amount <= 0 {
		http.Error(w, "Error of amount", http.StatusBadRequest)
		return
	}

	err := app.wallets.UpdateBalanceWithRetry(request.WalletID, int(request.Amount), request.OperationType, 5)
	if err != nil {
		app.serverError(w, err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (app *application) showBalance(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	fmt.Println(id)
	if len([]rune(id)) == 0 {
		app.notFound(w)
		return
	}
	s, err := app.wallets.GetBalance(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
			fmt.Println("Не найдено")
		}
		return
	}
	json.NewEncoder(w).Encode(map[string]int64{"balance": s})
}
