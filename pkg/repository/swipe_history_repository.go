// repository/SwipeHistory_repository.go
package repository

import (
	"errors"
	"log"
	"time"

	"github.com/metabbe3/knoxsdating/pkg/helpers"
	"github.com/metabbe3/knoxsdating/pkg/models"
)

type SwipeHistoryRepository interface {
	SaveSwipe(swipe *models.SwipeHistory, PremiumStatus interface{}) (bool, error)
	GetMatches(userID int, matchType string) ([]models.User, error)
	RedoSwipe(userID int) (*models.SwipeHistory, []models.User, error)
}

type swipeHistoryRepository struct {
	db    helpers.DatabaseHandler
	redis helpers.RedisHandler
}

func NewSwipeHistoryRepository(db helpers.DatabaseHandler, redis helpers.RedisHandler) SwipeHistoryRepository {
	return &swipeHistoryRepository{db: db, redis: redis}
}

// SaveSwipe saves the swipe history entry to the database and returns whether it is matched or not
func (r *swipeHistoryRepository) SaveSwipe(swipe *models.SwipeHistory, PremiumStatus interface{}) (bool, error) {
	// Check premium status and set the maximum allowed swipes
	maxSwipes := 10
	log.Println(PremiumStatus)
	if PremiumStatus == "Premium" {
		// If the user is premium, allow unlimited swipes
		maxSwipes = -1
	}

	// Fetch the user's total swipes for the day
	var totalSwipes int64
	r.db.Model(&models.SwipeHistory{}).
		Where(`"SwiperUserID" = ? AND DATE("Timestamp") = ?`, swipe.SwiperUserID, time.Now().UTC().Format("2006-01-02")).
		Count(&totalSwipes)

	// Check if the user has exceeded the maximum allowed swipes
	if maxSwipes != -1 && totalSwipes >= int64(maxSwipes) {
		return false, errors.New("maximum swipes exceeded for the day")
	}

	// Check if there is a match (opposite swipe direction from the swiped user)
	oppositeSwipe := models.SwipeHistory{}
	result := r.db.Where(
		`"SwiperUserID" = ? AND "SwipedUserID" = ? AND "SwipeDirection" = ?`,
		swipe.SwipedUserID, swipe.SwiperUserID, "right",
	).First(&oppositeSwipe)
	log.Print(result.RowsAffected)
	// If there is an opposite swipe, it's a match
	if result.RowsAffected > 0 && swipe.SwipeDirection == "right" {
		now := time.Now()
		// Update IsMatched in both entries
		swipe.IsMatched = true
		oppositeSwipe.IsMatched = true
		oppositeSwipe.Timestamp = now
		swipe.Timestamp = now

		// Reset RedoCount for both entries when it's a match
		swipe.RedoCount = 0

		// Save the updated entries
		r.db.Save(swipe)
		r.db.Save(&oppositeSwipe)

		return true, nil
	}

	// Create the swipe history entry if it's not a redo
	if swipe.RedoCount == 0 {
		r.db.Create(swipe)
	}

	return false, nil
}

func (r *swipeHistoryRepository) GetMatches(userID int, matchType string) ([]models.User, error) {
	var matches []models.User

	switch matchType {
	case "all":
		// Retrieve all matches where the current user has a match
		result := r.db.Table("SwipeHistory").
			Select(`DISTINCT ON ("SwipeHistory"."SwipedUserID") "SwipeHistory"."SwipedUserID", "User".*`).
			Joins(`JOIN "User" ON "SwipeHistory"."SwipedUserID" = "User"."UserID"`).
			Where(`"SwipeHistory"."SwiperUserID" = ? AND "SwipeHistory"."IsMatched" = true`, userID).
			Scan(&matches)
		if result.Error != nil {
			return nil, result.Error
		}

	case "liked":
		// Retrieve matches where the current user was liked by others
		result := r.db.Table("SwipeHistory").
			Select(`"User".*`).
			Joins(`JOIN "User" ON "SwipeHistory"."SwiperUserID" = "User"."UserID"`).
			Where(`"SwipeHistory"."SwipedUserID" = ? AND "SwipeHistory"."SwipeDirection" = 'right' AND "SwipeHistory"."IsMatched" = false`, userID).
			Scan(&matches)
		if result.Error != nil {
			return nil, result.Error
		}

	default:
		return nil, errors.New("invalid match type")
	}

	return matches, nil
}

// RedoSwipe allows a user to redo a swipe in the database
func (r *swipeHistoryRepository) RedoSwipe(userID int) (*models.SwipeHistory, []models.User, error) {
	// Retrieve the latest swipe entry with RedoCount > 0
	var originalSwipe models.SwipeHistory
	result := r.db.Where(`"SwiperUserID" = ?`, userID).Order(`"Timestamp" DESC`).First(&originalSwipe)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	if originalSwipe.RedoCount == 0 {
		// Perform the redo logic (e.g., go back to the previous profile)
		// ...

		// Increment the RedoCount
		originalSwipe.RedoCount++

		// Save the updated swipe entry back to the database
		result = r.db.Save(&originalSwipe)
		if result.Error != nil {
			return nil, nil, result.Error
		}

		// Retrieve the profiles and User based on the userID being swiped
		var profiles []models.User
		result = r.db.Table("SwipeHistory").
			Select(`DISTINCT ON ("SwipeHistory"."SwipedUserID") "SwipeHistory"."SwipedUserID", "User".*`).
			Joins(`JOIN "User" ON "SwipeHistory"."SwipedUserID" = "User"."UserID"`).
			Where(`"SwipeHistory"."SwiperUserID" = ? AND "SwipeHistory"."IsMatched" = true`, userID).
			Scan(&profiles)
		if result.Error != nil {
			return nil, nil, result.Error
		}

		return &originalSwipe, profiles, nil
	}

	return nil, nil, errors.New("no more redos available for this profile")
}

// NewUserRepositoryWithGormDBAndRedis creates a new ProfileRepository with GormDB and Redis
func NewSwipeHistoryRepositoryWithGormDBAndRedis(db helpers.DatabaseHandler, redis helpers.RedisHandler) SwipeHistoryRepository {
	return NewSwipeHistoryRepository(db, redis)
}
