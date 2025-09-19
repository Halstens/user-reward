package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

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
		return
	}
	var request struct {
		ReferrerID int `json:"referrer_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//Валидация
	if request.ReferrerID <= 0 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if id == request.ReferrerID {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	referrer, err := app.rewards.GetUserById(request.ReferrerID)
	if err != nil {
		if err == sql.ErrNoRows {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, err)
		}
		return
	}

	_, erro := app.rewards.GetUserById(id)
	if erro != nil {
		if erro == sql.ErrNoRows {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, erro)
		}
		return
	}

	alreadyReferred, err := app.isReferralExists(r.Context(), id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if alreadyReferred {
		app.clientError(w, http.StatusConflict)
		return
	}

	referralReward := 100

	err = app.processReferral(r.Context(), id, request.ReferrerID, referralReward)
	if err != nil {
		app.serverError(w, err)
		return
	}

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

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	user, err := app.rewards.GetUserByUsername(r.Context(), request.Username)
	if err != nil {
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	if request.Password != user.PasswordHash {
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = user.ID
	claims["username"] = request.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(app.jwt))
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":   tokenString,
		"message": "Login successful",
	})

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
