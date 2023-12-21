// handlers/swipe_history_handler.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/metabbe3/knoxsdating/pkg/helpers"
	"github.com/metabbe3/knoxsdating/pkg/models"
	"github.com/metabbe3/knoxsdating/pkg/repository"
)

type SwipeHistoryHandler struct {
	swipeHistoryRepo repository.SwipeHistoryRepository
	redisHelper      *helpers.RedisHelper
}

// NewLocationHandlers creates a new instance of LocationHandlers
func NewSwipeHistoryHandlers(swipeHistoryRepo repository.SwipeHistoryRepository, redisHelper *helpers.RedisHelper) *SwipeHistoryHandler {
	return &SwipeHistoryHandler{
		swipeHistoryRepo: swipeHistoryRepo,
		redisHelper:      redisHelper,
	}
}

// SaveSwipe handles the saving of swipe history entry
func (h *SwipeHistoryHandler) SaveSwipe(w http.ResponseWriter, r *http.Request) {
	var swipe models.SwipeHistory

	// Decode JSON payload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&swipe); err != nil {
		helpers.SendJSONResponse(w, http.StatusBadRequest, helpers.GenerateResponse(false, http.StatusBadRequest, "Invalid input", nil, err.Error()))
		return
	}
	defer r.Body.Close()

	tokenString := r.Header.Get("Authorization")
	token, err := helpers.ValidateToken(tokenString)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusUnauthorized, helpers.GenerateResponse(false, http.StatusUnauthorized, "Invalid token", nil, err.Error()))
		return
	}

	// Extract user ID from token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error parsing token claims", nil, ""))
		return
	}

	userID, ok := claims[helpers.UserIDKey].(float64)
	if !ok {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error parsing user ID from token claims", nil, ""))
		return
	}

	premiumStatus, ok := claims[helpers.PremiumStatusKey]
	if !ok {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error parsing user Premium Status from token claims", nil, ""))
		return
	}

	swipe.SwiperUserID = int(userID)

	isMatched, err := h.swipeHistoryRepo.SaveSwipe(&swipe, premiumStatus)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Failed to save swipe history", nil, err.Error()))
		return
	}

	// Update or create swipe history data in Redis
	err = h.updateOrCreateSwipeHistoryInRedis(int(userID), swipe)
	if err != nil {
		log.Printf("Error updating/creating swipe history in Redis: %v", err)
		// Handle the error as needed (e.g., log, but don't affect the HTTP response)
	}

	// Determine the match status based on IsMatched field
	matchStatus := "Not Matched"
	if isMatched {
		matchStatus = "Matched"

	}

	swipe.IsMatched = isMatched
	swipe.RedoCount = 0
	swipe.Timestamp = time.Now()

	// Use helpers.SendJSONResponse for the response
	helpers.SendJSONResponse(w, http.StatusOK, helpers.GenerateResponse(true, http.StatusOK, "Swipe history saved successfully", map[string]interface{}{
		"swipe":       swipe,
		"matchStatus": matchStatus,
	}, nil))
}

// GetMatches handles retrieving matches for a user
func (h *SwipeHistoryHandler) GetMatches(w http.ResponseWriter, r *http.Request) {
	// Decode JSON payload
	var requestBody struct {
		MatchType string `json:"matchType"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		helpers.SendJSONResponse(w, http.StatusBadRequest, helpers.GenerateResponse(false, http.StatusBadRequest, "Invalid input", nil, err.Error()))
		return
	}
	defer r.Body.Close()

	tokenString := r.Header.Get("Authorization")
	token, err := helpers.ValidateToken(tokenString)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusUnauthorized, helpers.GenerateResponse(false, http.StatusUnauthorized, "Invalid token", nil, err.Error()))
		return
	}

	// Extract user ID from token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error parsing token claims", nil, ""))
		return
	}

	userID, ok := claims[helpers.UserIDKey].(float64)
	if !ok {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error parsing user ID from token claims", nil, ""))
		return
	}

	// Call the repository method to get matches
	matches, err := h.swipeHistoryRepo.GetMatches(int(userID), requestBody.MatchType)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Failed to retrieve matches", nil, err.Error()))
		return
	}

	// Use helpers.SendJSONResponse for the response
	helpers.SendJSONResponse(w, http.StatusOK, helpers.GenerateResponse(true, http.StatusOK, "Matches retrieved successfully", matches, nil))
}

// RedoSwipe handles redoing a swipe
func (h *SwipeHistoryHandler) RedoSwipe(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from JWT token
	tokenString := r.Header.Get("Authorization")
	token, err := helpers.ValidateToken(tokenString)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusUnauthorized, helpers.GenerateResponse(false, http.StatusUnauthorized, "Invalid token", nil, err.Error()))
		return
	}

	// Extract user ID from token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error parsing token claims", nil, ""))
		return
	}

	userID, ok := claims[helpers.UserIDKey].(float64)
	if !ok {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error parsing user ID from token claims", nil, ""))
		return
	}

	// Retrieve redo result from Redis cache
	var redoResult string
	err = h.redisHelper.Get("redo:"+strconv.Itoa(int(userID)), &redoResult)
	if err != nil {
		// Handle the error, e.g., log it
		log.Printf("Error retrieving redo result from Redis: %v", err)
		redoResult = "" // Reset redoResult to proceed with the database operation
	}

	if redoResult == "error" {
		// If redo is not allowed, you may choose to log the event or return a specific response
		log.Println("Redo not allowed from Redis. Proceeding with database operation.")
	}

	// Update Redis cache with redo status
	err = h.redisHelper.Set("redo:"+strconv.Itoa(int(userID)), "error", 0)
	if err != nil {
		// Handle the error, e.g., log it
		log.Printf("Error updating Redis with redo status: %v", err)
		// Note: You may choose to continue with the operation or return an error response here
	}

	// Call the repository method to perform the redo swipe
	originalSwipe, profiles, err := h.swipeHistoryRepo.RedoSwipe(int(userID))
	if err != nil {
		// Handle the error, e.g., log it
		log.Printf("Error redoing swipe: %v", err)

		// If Redis did not have the result, you may want to attempt a database operation here
		if redoResult == "" {
			originalSwipe, profiles, err = h.swipeHistoryRepo.RedoSwipe(int(userID))
			if err != nil {
				// Handle the error from the database operation
				log.Printf("Error redoing swipe from DB: %v", err)
				helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Failed to redo swipe", nil, err.Error()))
				return
			}
		} else {
			helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Failed to redo swipe", nil, err.Error()))
			return
		}
	}

	// Use helpers.SendJSONResponse for the response
	responsePayload := map[string]interface{}{
		"originalSwipe": originalSwipe,
		"profiles":      profiles,
	}

	helpers.SendJSONResponse(w, http.StatusOK, helpers.GenerateResponse(true, http.StatusOK, "Redo swipe successful", responsePayload, nil))
}

// updateOrCreateSwipeHistoryInRedis updates or creates swipe history data in Redis
func (h *SwipeHistoryHandler) updateOrCreateSwipeHistoryInRedis(userID int, swipe models.SwipeHistory) error {
	// Construct the key based on your requirements
	key := "swipe_history:" + strconv.Itoa(userID)

	// Use the RedisHandler Set method to set the swipe history data in Redis with an expiration time (e.g., 24 hours)
	err := h.redisHelper.Set(key, swipe, time.Hour*24)
	return err
}
