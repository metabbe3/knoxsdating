// location_handlers.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/metabbe3/knoxsdating/pkg/helpers"
	"github.com/metabbe3/knoxsdating/pkg/models"
	"github.com/metabbe3/knoxsdating/pkg/repository"
)

type LocationHandlers struct {
	locationRepo repository.LocationRepository
	userRepo     repository.UserRepository
	redisHelper  *helpers.RedisHelper
}

// NewLocationHandlers creates a new instance of LocationHandlers
func NewLocationHandlers(locationRepo repository.LocationRepository, userRepo repository.UserRepository, redisHelper *helpers.RedisHelper) *LocationHandlers {
	return &LocationHandlers{
		locationRepo: locationRepo,
		userRepo:     userRepo,
		redisHelper:  redisHelper,
	}
}

func (h *LocationHandlers) CreateLocationHistory(w http.ResponseWriter, r *http.Request) {
	var location models.LocationHistory

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&location); err != nil {
		helpers.SendJSONResponse(w, http.StatusBadRequest, helpers.GenerateResponse(false, http.StatusBadRequest, "Invalid request payload", nil, err.Error()))
		return
	}

	defer r.Body.Close()

	// Validate JWT token
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

	// Set the UserID field of the location struct
	location.UserID = int(userID)

	// Fetch the user data from the user repository
	user, err := h.userRepo.GetUserByID(int(userID))
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error fetching user data", nil, err.Error()))
		return
	}

	// Check if the user is premium
	isPremium := user.PremiumStatus == "Premium"

	// If the user is not premium, set the Timestamp to the current time
	if !isPremium {
		location.Timestamp = time.Now()
	}

	// Call the repository method without worrying about JWT token validation
	err = h.locationRepo.CreateLocationHistory(&location, isPremium)
	if err != nil {
		// If the user is premium, allow unlimited entries; otherwise, check for duplicate key violation
		if isPremium {
			helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error creating location history", nil, err.Error()))
		} else if strings.Contains(err.Error(), "location history already created for the user today") {
			helpers.SendJSONResponse(w, http.StatusConflict, helpers.GenerateResponse(false, http.StatusConflict, "Location history already exists for today", nil, err.Error()))
		} else {
			helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error creating location history", nil, err.Error()))
		}
		return
	}

	// Update or create location data in Redis
	err = h.updateOrCreateLocationInRedis(userID, location)
	if err != nil {
		log.Printf("Error updating/creating location in Redis: %v", err)
	}

	helpers.SendJSONResponse(w, http.StatusCreated, helpers.GenerateResponse(true, http.StatusCreated, "Location history created successfully", location, nil))
}

// Helper function to update or create location data in Redis
func (h *LocationHandlers) updateOrCreateLocationInRedis(userID float64, location models.LocationHistory) error {
	// Delete existing location data from Redis
	err := h.redisHelper.Delete("location:" + strconv.Itoa(int(userID)))
	if err != nil {
		return err
	}

	// Set location data in Redis
	err = h.redisHelper.Set("location:"+strconv.Itoa(int(userID)), location, time.Hour*24)
	return err
}

// GetNearbyLocations handles the request to get nearby locations with pagination
func (h *LocationHandlers) GetNearbyLocations(w http.ResponseWriter, r *http.Request) {
	// Validate JWT token
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

	// Decode JSON payload
	var requestPayload struct {
		MaxDistance float64 `json:"maxDistance"`
		Page        int     `json:"page"`
		PageSize    int     `json:"pageSize"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestPayload); err != nil {
		helpers.SendJSONResponse(w, http.StatusBadRequest, helpers.GenerateResponse(false, http.StatusBadRequest, "Invalid request payload", nil, err.Error()))
		return
	}

	// Fetch user's location from Redis
	var userLocation models.LocationHistory
	if err := h.redisHelper.Get("location:"+strconv.Itoa(int(userID)), &userLocation); err != nil {
		// If not found in Redis, fetch from the database
		locationHistory, err := h.locationRepo.GetLocationHistoryByUserID(int(userID))
		if err != nil {
			helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error fetching user's location", nil, err.Error()))
			return
		}

		// Check if the user's location is available in the database
		if len(locationHistory) == 0 {
			helpers.SendJSONResponse(w, http.StatusNotFound, helpers.GenerateResponse(false, http.StatusNotFound, "User's location not found", nil, ""))
			return
		}

		// Use the latest location from the database
		userLocation = locationHistory[0]

		// Store the user's location in Redis for future use
		err = h.redisHelper.Set("location:"+strconv.Itoa(int(userID)), userLocation, time.Hour*24)
		if err != nil {
			log.Printf("Error storing user's location in Redis: %v", err)
		}
	}

	// Fetch nearby locations
	nearbyLocations, err := h.locationRepo.GetNearbyLocations(int(userID), requestPayload.MaxDistance, requestPayload.Page, requestPayload.PageSize)

	if err != nil {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error fetching nearby locations", nil, err.Error()))
		return
	}

	// Send the response
	helpers.SendJSONResponse(w, http.StatusOK, helpers.GenerateResponse(true, http.StatusOK, "Nearby locations fetched successfully", nearbyLocations, nil))
}
