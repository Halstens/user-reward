package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/user-reward/internal/models"
)

func (app *application) ShowUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow: ", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		app.notFound(w)
		fmt.Println("Нет ид")
		return
	}
	info, err := app.rewards.GetUserById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
			fmt.Println("Не найдено", id)
		}
		return
	}
	json.NewEncoder(w).Encode(info)

}

func (app *application) ShowTopUserByBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow: ", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	topList, err := app.rewards.GetTopList()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
			fmt.Println("Не найдено")
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topList)

}

func (app *application) CompletedTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow: ", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		app.notFound(w)
		fmt.Println("Нет ид")
		return
	}
	var request struct {
		TaskType string `json:"task_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	validTasks := map[string]int{
		"subscribe_telegram": 50,
		"follow_twitter":     30,
		"referral_signup":    100,
	}
	reward, exists := validTasks[request.TaskType]
	if !exists {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	user, err := app.rewards.GetUserById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	alreadyCompleted, err := app.rewards.IsTaskCompleted(r.Context(), id, request.TaskType)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if alreadyCompleted {
		app.clientError(w, http.StatusConflict)
		return
	}

	err = app.rewards.UpdateUserBalance(id, reward, models.OperationType(request.TaskType))
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.rewards.MarkTaskCompleted(r.Context(), id, request.TaskType)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"message":     "Task completed successfully",
		"reward":      reward,
		"new_balance": user.Balance + reward,
		"task_type":   request.TaskType,
	})

}

func (app *application) isReferralExists(ctx context.Context, refereeID int) (bool, error) {
	return app.rewards.IsReferralExists(ctx, refereeID)
}

func (app *application) processReferral(ctx context.Context, refereeID, referrerID, reward int) error {
	return app.rewards.ProcessReferral(ctx, referrerID, refereeID, reward)
}

func (app *application) AddRefferer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow: ", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		app.notFound(w)
		fmt.Println("Нет ид")
		return
	}
	var request struct {
		ReferrerID int `json:"referrer_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		app.clientError(w, http.StatusBadRequest)
		fmt.Println("1", request.ReferrerID)
		return
	}

	// Валидация
	if request.ReferrerID <= 0 {
		app.clientError(w, http.StatusBadRequest)
		fmt.Println("2", request.ReferrerID)
		return
	}

	// Нельзя указать самого себя как реферера
	if id == request.ReferrerID {
		app.clientError(w, http.StatusBadRequest)
		fmt.Println("3", request.ReferrerID)
		return
	}

	// Проверяем, существует ли реферер
	referrer, err := app.rewards.GetUserById(request.ReferrerID)
	if err != nil {
		if err == sql.ErrNoRows {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Проверяем, существует ли реферал
	_, erro := app.rewards.GetUserById(id)
	if erro != nil {
		if erro == sql.ErrNoRows {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, erro)
		}
		return
	}

	// Проверяем, не добавлен ли уже реферер для этого пользователя
	alreadyReferred, err := app.isReferralExists(r.Context(), id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if alreadyReferred {
		app.clientError(w, http.StatusConflict)
		return
	}

	// Награда за реферала
	referralReward := 100

	// Выполняем в транзакции
	err = app.processReferral(r.Context(), id, request.ReferrerID, referralReward)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       true,
		"message":       "Referrer added successfully",
		"referrer_id":   request.ReferrerID,
		"referrer_name": referrer.Name,
		"reward":        referralReward,
		"new_balance":   referrer.Balance + referralReward,
	})
}
