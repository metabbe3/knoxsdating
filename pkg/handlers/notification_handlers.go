// Import the required packages
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/metabbe3/knoxsdating/pkg/helpers"
	"github.com/metabbe3/knoxsdating/pkg/models"
	"github.com/metabbe3/knoxsdating/pkg/repository"
)

type NotificationHandlers struct {
	notificationRepo repository.NotificationRepository
	redisHelper      *helpers.RedisHelper
}

func NewNotificationHandlers(repo repository.NotificationRepository, redisHelper *helpers.RedisHelper) *NotificationHandlers {
	return &NotificationHandlers{
		notificationRepo: repo,
		redisHelper:      redisHelper,
	}
}

func (h *NotificationHandlers) CreateNotification(w http.ResponseWriter, r *http.Request) {
	var notification models.Notification

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&notification); err != nil {
		helpers.SendJSONResponse(w, http.StatusBadRequest, helpers.GenerateResponse(false, http.StatusBadRequest, "Invalid request payload", nil, err.Error()))
		return
	}

	defer r.Body.Close()

	// Validate JWT token
	tokenString := r.Header.Get("Authorization")
	log.Print(tokenString)
	_, err := helpers.ValidateToken(tokenString)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusUnauthorized, helpers.GenerateResponse(false, http.StatusUnauthorized, "Invalid token", nil, err.Error()))
		return
	}

	err = h.notificationRepo.CreateNotification(&notification)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error creating notification", nil, err.Error()))
		log.Printf("Error creating notification: %v", err)
		return
	}

	// Cache notification data in Redis
	err = h.redisHelper.Set("notification:"+strconv.Itoa(notification.NotificationID), notification, time.Hour*24)
	if err != nil {
		log.Printf("Error caching notification data in Redis: %v", err)
		// Handle error, e.g., log it
	}

	helpers.SendJSONResponse(w, http.StatusCreated, helpers.GenerateResponse(true, http.StatusCreated, "Notification created successfully", notification, nil))
}

func (h *NotificationHandlers) GetNotificationByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusBadRequest, helpers.GenerateResponse(false, http.StatusBadRequest, "Invalid notification ID", nil, err.Error()))
		return
	}

	// Validate JWT token
	tokenString := r.Header.Get("Authorization")
	log.Print(tokenString)
	_, err = helpers.ValidateToken(tokenString)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusUnauthorized, helpers.GenerateResponse(false, http.StatusUnauthorized, "Invalid token", nil, err.Error()))
		return
	}

	// Try to get notification data from Redis first
	var cachedNotification models.Notification
	err = h.redisHelper.Get("notification:"+strconv.Itoa(id), &cachedNotification)
	if err == nil {
		// Use cached notification data if available
		helpers.SendJSONResponse(w, http.StatusOK, helpers.GenerateResponse(true, http.StatusOK, "Notification retrieved successfully", cachedNotification, nil))
		return
	}

	// If not found in Redis, fetch from the database
	notification, err := h.notificationRepo.GetNotificationByID(id)
	if err != nil {
		log.Printf("Error fetching notification: %v", err)
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error fetching notification", nil, err.Error()))
		return
	}

	// Cache fetched notification data in Redis
	err = h.redisHelper.Set("notification:"+strconv.Itoa(id), notification, time.Hour*24)
	if err != nil {
		log.Printf("Error caching notification data in Redis: %v", err)
		// Handle error, e.g., log it
	}

	helpers.SendJSONResponse(w, http.StatusOK, helpers.GenerateResponse(true, http.StatusOK, "Notification retrieved successfully", notification, nil))
}

func (h *NotificationHandlers) UpdateNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusBadRequest, helpers.GenerateResponse(false, http.StatusBadRequest, "Invalid notification ID", nil, err.Error()))
		return
	}

	// Validate JWT token
	tokenString := r.Header.Get("Authorization")
	log.Print(tokenString)
	_, err = helpers.ValidateToken(tokenString)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusUnauthorized, helpers.GenerateResponse(false, http.StatusUnauthorized, "Invalid token", nil, err.Error()))
		return
	}

	var notification models.Notification

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&notification); err != nil {
		helpers.SendJSONResponse(w, http.StatusBadRequest, helpers.GenerateResponse(false, http.StatusBadRequest, "Invalid request payload", nil, err.Error()))
		return
	}

	defer r.Body.Close()

	notification.NotificationID = id
	err = h.notificationRepo.UpdateNotification(&notification)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error updating notification", nil, err.Error()))
		log.Printf("Error updating notification: %v", err)
		return
	}

	// Update cached notification data in Redis
	err = h.redisHelper.Set("notification:"+strconv.Itoa(id), notification, time.Hour*24)
	if err != nil {
		log.Printf("Error updating cached notification data in Redis: %v", err)
		// Handle error, e.g., log it
	}

	helpers.SendJSONResponse(w, http.StatusOK, helpers.GenerateResponse(true, http.StatusOK, "Notification updated successfully", notification, nil))
}

func (h *NotificationHandlers) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusBadRequest, helpers.GenerateResponse(false, http.StatusBadRequest, "Invalid notification ID", nil, err.Error()))
		return
	}

	// Validate JWT token
	tokenString := r.Header.Get("Authorization")
	log.Print(tokenString)
	_, err = helpers.ValidateToken(tokenString)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusUnauthorized, helpers.GenerateResponse(false, http.StatusUnauthorized, "Invalid token", nil, err.Error()))
		return
	}

	// Get notification data from Redis before deletion
	var cachedNotification models.Notification
	err = h.redisHelper.Get("notification:"+strconv.Itoa(id), &cachedNotification)
	if err != nil {
		log.Printf("Error fetching cached notification data from Redis: %v", err)
		// Handle error, e.g., log it
	}

	// Delete notification data from Redis
	err = h.redisHelper.Delete("notification:" + strconv.Itoa(id))
	if err != nil {
		log.Printf("Error deleting cached notification data in Redis: %v", err)
		// Handle error, e.g., log it
	}

	// Delete notification from the database
	err = h.notificationRepo.DeleteNotification(&cachedNotification)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error deleting notification", nil, err.Error()))
		log.Printf("Error deleting notification: %v", err)
		return
	}

	helpers.SendJSONResponse(w, http.StatusOK, helpers.GenerateResponse(true, http.StatusOK, "Notification deleted successfully", nil, nil))
}
