package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/metabbe3/knoxsdating/pkg/helpers"
	"github.com/metabbe3/knoxsdating/pkg/models"
	"github.com/metabbe3/knoxsdating/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserHandlers struct {
	userRepo    repository.UserRepository
	redisHelper *helpers.RedisHelper
}

func NewUserHandlers(userRepo repository.UserRepository, redisHelper *helpers.RedisHelper) *UserHandlers {
	return &UserHandlers{
		userRepo:    userRepo,
		redisHelper: redisHelper,
	}
}

func (h *UserHandlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		helpers.SendJSONResponse(w, http.StatusBadRequest, helpers.GenerateResponse(false, http.StatusBadRequest, "Invalid request payload", nil, err.Error()))
		return
	}

	defer r.Body.Close()

	// Validate input
	if err := validateUserInput(&user); err != nil {
		helpers.SendJSONResponse(w, http.StatusBadRequest, helpers.GenerateResponse(false, http.StatusBadRequest, "Invalid user input", nil, err.Error()))
		return
	}

	// Check if the email already exists
	existingUser, _ := h.userRepo.GetUserByEmail(user.Email)
	if existingUser != nil {
		helpers.SendJSONResponse(w, http.StatusConflict, helpers.GenerateResponse(false, http.StatusConflict, "Email already exists", nil, nil))
		return
	}

	// Validate password
	if len(user.Password) < 8 {
		helpers.SendJSONResponse(w, http.StatusBadRequest, helpers.GenerateResponse(false, http.StatusBadRequest, "Password must be at least 8 characters long", nil, nil))
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error hashing password", nil, err.Error()))
		return
	}
	user.Password = string(hashedPassword)

	// Create user
	err = h.userRepo.CreateUser(&user)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error creating user", nil, err.Error()))
		return
	}

	// Set user data in Redis
	err = h.redisHelper.Set("user:"+user.Email, user, time.Hour*24)
	if err != nil {
		// Handle error, e.g., log it
		log.Printf("Error setting data in Redis: %v", err)
	}

	helpers.SendJSONResponse(w, http.StatusCreated, helpers.GenerateResponse(true, http.StatusCreated, "User created successfully", user, nil))
}

func (h *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&credentials); err != nil {
		helpers.SendJSONResponse(w, http.StatusBadRequest, helpers.GenerateResponse(false, http.StatusBadRequest, "Invalid request payload", nil, err.Error()))
		return
	}

	defer r.Body.Close()

	// Validate input
	if err := validateLoginInput(&credentials); err != nil {
		helpers.SendJSONResponse(w, http.StatusBadRequest, helpers.GenerateResponse(false, http.StatusBadRequest, "Invalid login input", nil, err.Error()))
		return
	}

	// Retrieve user data from Redis
	var cachedUser models.User
	err := h.redisHelper.Get("user:"+credentials.Email, &cachedUser)
	if err != nil {
		// If user data is not in Redis, fetch it from the database
		user, err := h.userRepo.GetUserByEmail(credentials.Email)
		if err != nil {
			helpers.SendJSONResponse(w, http.StatusUnauthorized, helpers.GenerateResponse(false, http.StatusUnauthorized, "Invalid username or password", nil, err.Error()))
			return
		}

		// Cache hashed password in Redis
		err = h.redisHelper.Set("user:"+credentials.Email, user, time.Hour*24)
		if err != nil {
			// Handle error, e.g., log it
			log.Printf("Error caching data in Redis: %v", err)
		}

		// Use fetched user data
		cachedUser = *user
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(cachedUser.Password), []byte(credentials.Password))
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusUnauthorized, helpers.GenerateResponse(false, http.StatusUnauthorized, "Invalid username or password", nil, err.Error()))
		return
	}

	// Generate JWT token
	token, err := helpers.GenerateToken(cachedUser)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error generating token", nil, err.Error()))
		return
	}

	// Return success response with token
	response := models.LoginResponse{
		Token: token,
		User:  cachedUser, // Use the cached user data here
	}
	helpers.SendJSONResponse(w, http.StatusOK, helpers.GenerateResponse(true, http.StatusOK, "Login successful", response, nil))
}

func (h *UserHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Declare and initialize the secretKey variable
	secretKey := r.Header.Get("SecretKey")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		helpers.SendJSONResponse(w, http.StatusBadRequest, helpers.GenerateResponse(false, http.StatusBadRequest, "Invalid request payload", nil, err.Error()))
		return
	}

	defer r.Body.Close()

	// Validate JWT token
	tokenString := r.Header.Get("Authorization")
	log.Print(tokenString)
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

	// Fetch the existing user data
	existingUser, err := h.userRepo.GetUserByEmail(email)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error fetching user", nil, err.Error()))
		return
	}

	if secretKey == "adminSecretKey" {
		// Update the user data based on the input
		if user.Gender != "" {
			existingUser.Gender = user.Gender
		}
		if user.Company != "" {
			existingUser.Company = user.Company
		}
		if user.School != "" {
			existingUser.School = user.School
		}
		if user.JobTitle != "" {
			existingUser.JobTitle = user.JobTitle
		}
	} else {
		// Update the user data based on the input, excluding certain fields
		if user.Gender != "" {
			existingUser.Gender = user.Gender
		}
		if user.Company != "" {
			existingUser.Company = user.Company
		}
		if user.School != "" {
			existingUser.School = user.School
		}
		if user.JobTitle != "" {
			existingUser.JobTitle = user.JobTitle
		}

		// Check if the password is provided in the input
		if user.Password != "" {
			// Check if the password is already hashed
			if !strings.HasPrefix(user.Password, "$2a$") {
				// Hash the new password
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
				if err != nil {
					helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error hashing password", nil, err.Error()))
					return
				}
				existingUser.Password = string(hashedPassword)
			} else {
				// The password is already hashed
				existingUser.Password = user.Password
			}
		}
	}

	// Update the user by email
	err = h.userRepo.UpdateUser(existingUser)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error updating user", nil, err.Error()))
		return
	}

	// Delete user data from Redis on update
	err = h.redisHelper.Delete("user:" + existingUser.Email)
	if err != nil {
		log.Printf("Error deleting data in Redis: %v", err)
		// Handle error, e.g., log it
	}

	// Add back updated data to Redis
	err = h.redisHelper.Set("user:"+existingUser.Email, existingUser, time.Hour*24)
	if err != nil {
		log.Printf("Error setting data in Redis: %v", err)
		// Handle error, e.g., log it
	}

	helpers.SendJSONResponse(w, http.StatusOK, helpers.GenerateResponse(true, http.StatusOK, "User updated successfully", existingUser, nil))
}

func (h *UserHandlers) UpdatePremium(w http.ResponseWriter, r *http.Request) {
	var user models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
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

	// Fetch the existing user data
	existingUser, err := h.userRepo.GetUserByEmail(email)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error fetching user", nil, err.Error()))
		return
	}

	// Check if the user is already a premium user
	if existingUser.PremiumStatus == "Free" {
		// Set the PremiumStatus to "Premium"
		existingUser.PremiumStatus = "Premium"

		// Set the PremiumStartDate to the current date
		currentTime := time.Now()
		existingUser.PremiumStartDate = currentTime

		existingUser.PremiumEndDate = currentTime.AddDate(0, 1, 0)

	} else if existingUser.PremiumStatus == "Premium" {
		// Add 1 month to the current PremiumEndDate
		// Set the PremiumEndDate to 1 month from the current date or the token's end date + 1 month
		existingUser.PremiumEndDate = existingUser.PremiumEndDate.AddDate(0, 1, 0)
	}

	// Update the user by email
	err = h.userRepo.UpdateUser(existingUser)
	if err != nil {
		helpers.SendJSONResponse(w, http.StatusInternalServerError, helpers.GenerateResponse(false, http.StatusInternalServerError, "Error updating user", nil, err.Error()))
		return
	}

	// Delete user data from Redis on update
	err = h.redisHelper.Delete("user:" + existingUser.Email)
	if err != nil {
		log.Printf("Error delete  data in Redis: %v", err)
		// Handle error, e.g., log it
	}

	// Add back updated data to Redis
	err = h.redisHelper.Set("user:"+existingUser.Email, existingUser, time.Hour*24)
	if err != nil {
		log.Printf("Error set  data in Redis: %v", err)
		// Handle error, e.g., log it
	}

	helpers.SendJSONResponse(w, http.StatusOK, helpers.GenerateResponse(true, http.StatusOK, "User updated successfully", existingUser, nil))
}

func validateUserInput(user *models.User) error {
	// Validate email
	if !helpers.IsValidEmail(user.Email) {
		return helpers.ValidationError("Invalid email")
	}

	// Validate password
	if len(user.Password) < 8 {
		return helpers.ValidationError("Password must be at least 8 characters long")
	}

	// Validate username
	if len(user.Username) == 0 {
		return helpers.ValidationError("Username is required")
	}

	// Validate gender
	if len(user.Gender) > 0 && (user.Gender != "Male" && user.Gender != "Female") {
		return helpers.ValidationError("Invalid gender")
	}

	return nil
}

func validateLoginInput(credentials *models.Credentials) error {
	// Validate email
	if !helpers.IsValidEmail(credentials.Email) {
		return helpers.ValidationError("Invalid email")
	}

	// Validate password
	if len(credentials.Password) < 8 {
		return helpers.ValidationError("Password must be at least 8 characters long")
	}

	return nil
}
