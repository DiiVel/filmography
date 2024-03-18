package handlers

import (
	"context"
	"encoding/json"
	"filmography/internal/entities"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserService interface {
	CreateUser(ctx context.Context, user entities.UserEntity) error
	GetUsers(ctx context.Context) ([]entities.UserEntity, error)
	GetUser(ctx context.Context, id string) (entities.UserEntity, error)
	UpdateUser(ctx context.Context, id string, user entities.UserEntity) error
	DeleteUser(ctx context.Context, id string) error
}

// createUser создает нового юзера.
// @Summary Создает юзера.
// @Description Создает нового юзера на основе переданных данных.
// @Tags User
// @Accept json
// @Produce json
// @Param user body entities.UserEntity true "Данные юзера"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Ошибка при декодировании JSON"
// @Failure 500 {string} string "Ошибка при создании юзера"
// @Router /user [post]
func (handlers Handlers) createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed", r.Method)
		return
	}

	var user entities.UserEntity
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("decode request body failed")
		return
	}

	err = handlers.svc.CreateUser(context.Background(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating user: %v", err)
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("create user failed")
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User successfully created")
}

// getUsers возвращает список юзеров.
// @Summary Возвращает список юзеров
// @Description Возвращает список всех юзеров.
// @Tags User
// @Produce json
// @Success 200 {array} entities.UserEntity "Список юзеров"
// @Failure 500 {string} string "Ошибка при получении юзеров"
// @Router /user [get]
func (handlers Handlers) getUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed", r.Method)
		return
	}

	users, err := handlers.svc.GetUsers(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting users: %v", err)
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("get users failed")
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("encode users failed")
		return
	}
}

// getUser возвращает информацию о юзере по его ID.
// @Summary Возвращает информацию о юзере
// @Description Возвращает информацию о юзере по указанному ID.
// @Tags User
// @Param id query string true "ID юзера"
// @Produce json
// @Success 200 {object} entities.UserEntity "Информация о юзере"
// @Failure 500 {string} string "Ошибка при получении юзере"
// @Router /user/{id} [get]
func (handlers Handlers) getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed", r.Method)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing id parameter")
		return
	}

	user, err := handlers.svc.GetUser(context.Background(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting user: %v", err)
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("get user failed")
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("encode user failed")
		return
	}
}

// updateUser обновляет информацию о юзере.
// @Summary Обновляет информацию о юзере
// @Description Обновляет информацию о юзере с указанным ID на основе переданных данных.
// @Tags User
// @Param id query string true "ID юзера"
// @Accept json
// @Produce json
// @Param user body entities.UserEntity true "Данные юзера"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Ошибка при декодировании JSON"
// @Failure 500 {string} string "Ошибка при обновлении юзера"
// @Router /user/{id} [put]
func (handlers Handlers) updateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed", r.Method)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing id parameter")
		return
	}

	var user entities.UserEntity
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("decode request body failed")
		return
	}

	err = handlers.svc.UpdateUser(context.Background(), id, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating user: %v", err)
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("update user failed")
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User successfully updated")
	return
}

// deleteUser удаляет юзера по его ID.
// @Summary Удаляет юзера
// @Description Удаляет юзера с указанным ID.
// @Tags User
// @Param id query string true "ID юзера"
// @Success 200 {object} map[string]string
// @Failure 500 {string} string "Ошибка при удалении юзера"
// @Router /user/{id} [delete]
func (handlers Handlers) deleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed", r.Method)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing id parameter")
		return
	}

	err := handlers.svc.DeleteUser(context.Background(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting user: %v", err)
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("delete user failed")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User successfully deleted")
}
