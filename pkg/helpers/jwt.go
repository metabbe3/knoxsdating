package helpers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/metabbe3/knoxsdating/pkg/models"
)

const (
	// Define a secret key for signing and verifying the token
	secretKey = "your-secret-key"

	// Define constants for user data keys
	UserIDKey             = "userID"
	UsernameKey           = "username"
	EmailKey              = "email"
	VerificationStatusKey = "verificationStatus"
	VerificationBadgeKey  = "verificationBadge"
	PremiumStatusKey      = "premiumStatus"
	PremiumStartDateKey   = "premiumStartDate"
	PremiumEndDateKey     = "premiumEndDate"
	GenderKey             = "gender"
	CompanyKey            = "company"
	SchoolKey             = "school"
	JobTitleKey           = "jobTitle"
	VerifiedBadgeKey      = "verifiedBadge"
)

func GenerateToken(user models.User) (string, error) {
	// Create the claims for the token
	claims := jwt.MapClaims{
		UserIDKey:             user.UserID,
		UsernameKey:           user.Username,
		EmailKey:              user.Email,
		VerificationStatusKey: user.VerificationStatus,
		VerificationBadgeKey:  user.VerificationBadge,
		PremiumStatusKey:      user.PremiumStatus,
		PremiumStartDateKey:   user.PremiumStartDate,
		PremiumEndDateKey:     user.PremiumEndDate,
		GenderKey:             user.Gender,
		CompanyKey:            user.Company,
		SchoolKey:             user.School,
		JobTitleKey:           user.JobTitle,
		VerifiedBadgeKey:      user.VerifiedBadge,
		// Add more user data as needed
	}

	// Set the expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims["exp"] = expirationTime.Unix()

	// Create the token with the claims and sign it with the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return signedToken, nil
}

func ReadAndDecryptToken(tokenString string) (jwt.MapClaims, error) {
	// Remove the "bearer " prefix from the token string
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method and return the secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to extract claims from token")
	}

	return claims, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	// Remove the "bearer " prefix from the token string
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method and return the secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}
