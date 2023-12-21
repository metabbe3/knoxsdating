// profile_handlers.go
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

type ProfileHandlers struct {
	profileRepo repository.ProfileRepository
	redisHelper *helpers.RedisHelper
}

// NewProfileHandlers creates a new instance of ProfileHandlers
func NewProfileHandlers(profileRepo repository.ProfileRepository, redisHelper *helpers.RedisHelper) *ProfileHandlers {
	return &ProfileHandlers{
		profileRepo: profileRepo,
		redisHelper: redisHelper,
	}
}

func (h *ProfileHandlers) CreateProfile(w http.ResponseWriter, r *http.Request) {
	var profile models.Profile

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&profile); err != nil {
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

	// Get the user email from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error parsing token claims", nil, ""))
		return
	}
	email, ok := claims["email"].(string)
	if !ok {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error parsing email from token claims", nil, ""))
		return
	}
	log.Println("Before Redis Get", email)
	// Fetch the user data from Redis
	var user models.User
	err = h.redisHelper.Get("user:"+email, &user)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error fetching user from Redis", nil, err.Error()))
		return
	}
	log.Println("After Redis Get", user.UserID, user.Email)

	// Get the existing profile of the user

	existingProfile, err := h.profileRepo.GetProfileByUserID(user.UserID)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error fetching existing profile", nil, err.Error()))
		return
	}

	// If the user already has a profile, update it; otherwise, create a new profile
	if existingProfile != nil {
		// Update the existing profile with user input
		updateProfile(existingProfile, &profile)

		// Update the profile in the database
		err = h.profileRepo.UpdateProfile(existingProfile)
		if err != nil {
			helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error updating profile", nil, err.Error()))
			return
		}

		// Set updated profile data in Redis
		err = h.redisHelper.Set("profile:"+strconv.Itoa(existingProfile.UserID), existingProfile, time.Hour*24)
		if err != nil {
			log.Printf("Error setting profile data in Redis: %v", err)
		}

		helpers.SendJSONResponse(w, http.StatusOK, helpers.GenerateResponse(true, http.StatusOK, "Profile updated successfully", existingProfile, nil))
	} else {
		// Associate the profile with the user
		profile.UserID = user.UserID

		// Create a new profile for the user
		err = h.profileRepo.CreateProfile(&profile)
		if err != nil {
			helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error creating profile", nil, err.Error()))
			return
		}

		// Set profile data in Redis
		err = h.redisHelper.Set("profile:"+strconv.Itoa(profile.UserID), profile, time.Hour*24)
		if err != nil {
			log.Printf("Error setting profile data in Redis: %v", err)
		}

		helpers.SendJSONResponse(w, http.StatusCreated, helpers.GenerateResponse(true, http.StatusCreated, "Profile created successfully", profile, nil))
	}
}

func (h *ProfileHandlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Validate JWT token
	tokenString := r.Header.Get("Authorization")
	token, err := helpers.ValidateToken(tokenString)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusUnauthorized, helpers.GenerateResponse(false, http.StatusUnauthorized, "Invalid token", nil, err.Error()))
		return
	}

	// Get the user email from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error parsing token claims", nil, ""))
		return
	}
	email, ok := claims["email"].(string)
	if !ok {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error parsing email from token claims", nil, ""))
		return
	}

	// Fetch the user data from Redis
	var user models.User
	err = h.redisHelper.Get("user:"+email, &user)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error fetching user from Redis", nil, err.Error()))
		return
	}

	// Check if the profile data is in Redis based on userID
	var cachedProfile models.Profile
	err = h.redisHelper.Get("profile:"+strconv.Itoa(user.UserID), &cachedProfile)
	if err == nil {
		helpers.SendJSONResponse(w, http.StatusOK, helpers.GenerateResponse(true, http.StatusOK, "Profile retrieved successfully", cachedProfile, nil))
		return
	}

	// Fetch the profile data from the database by userID
	profile, err := h.profileRepo.GetProfileByUserID(user.UserID)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error fetching profile", nil, err.Error()))
		return
	}

	// Set profile data in Redis
	err = h.redisHelper.Set("profile:"+strconv.Itoa(profile.UserID), profile, time.Hour*24)
	if err != nil {
		log.Printf("Error setting profile data in Redis: %v", err)
	}

	helpers.SendJSONResponse(w, http.StatusOK, helpers.GenerateResponse(true, http.StatusOK, "Profile retrieved successfully", profile, nil))
}

func updateProfile(profile *models.Profile, userInput *models.Profile) {
	// Check if the user input contains the "Photos" field
	if userInput.Photos != "" {
		profile.Photos = userInput.Photos
	}

	// Check if the user input contains the "AboutMe" field
	if userInput.AboutMe != "" {
		profile.AboutMe = userInput.AboutMe
	}

	// Check if the user input contains the "Interests" field
	if userInput.Interests != "" {
		profile.Interests = userInput.Interests
	}

	// Check if the user input contains the "RelationshipGoals" field
	if userInput.RelationshipGoals != "" {
		profile.RelationshipGoals = userInput.RelationshipGoals
	}

	// Check if the user input contains the "Height" field
	if userInput.Height != 0 {
		profile.Height = userInput.Height
	}

	// Check if the user input contains the "Language" field
	if userInput.Language != "" {
		profile.Language = userInput.Language
	}

	// Check if the user input contains the "ZodiacSign" field
	if userInput.ZodiacSign != "" {
		profile.ZodiacSign = userInput.ZodiacSign
	}

	// Check if the user input contains the "EducationDetails" field
	if userInput.EducationDetails != "" {
		profile.EducationDetails = userInput.EducationDetails
	}

	// Check if the user input contains the "SocialMediaAccounts" field
	if userInput.SocialMediaAccounts != "" {
		profile.SocialMediaAccounts = userInput.SocialMediaAccounts
	}
}

// Add other profile-related handlers as needed
