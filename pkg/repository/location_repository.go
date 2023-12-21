// location_repository.go
package repository

import (
	"errors"
	"log"
	"time"

	"github.com/metabbe3/knoxsdating/pkg/helpers"
	"github.com/metabbe3/knoxsdating/pkg/models"
)

type LocationRepository interface {
	CreateLocationHistory(location *models.LocationHistory, isPremium bool) error
	GetLocationHistoryByUserID(userID int) ([]models.LocationHistory, error)
	GetNearbyLocations(userID int, maxDistance float64, page, pageSize int) ([]LocationWithDistance, error)
}

type locationRepository struct {
	db          helpers.DatabaseHandler
	redis       helpers.RedisHandler
	profileRepo ProfileRepository
}

// NearbyLocation represents the structure of nearby locations
type NearbyLocation struct {
	LocationID int       `json:"locationID"`
	UserID     int       `json:"userID"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Timestamp  time.Time `json:"timestamp"`
	Distance   float64   `json:"distance"`
}

type LocationWithDistance struct {
	LocationID int             `json:"locationID"`
	UserID     int             `json:"userID"`
	Latitude   float64         `json:"latitude"`
	Longitude  float64         `json:"longitude"`
	Timestamp  time.Time       `json:"timestamp"`
	Distance   float64         `json:"distance"`
	Profile    *models.Profile `json:"profile,omitempty" gorm:"foreignKey:UserID"`
}

func NewLocationRepository(db helpers.DatabaseHandler, redis helpers.RedisHandler, profileRepo ProfileRepository) LocationRepository {
	return &locationRepository{
		db:          db,
		redis:       redis,
		profileRepo: profileRepo,
	}
}

func (r *locationRepository) CreateLocationHistory(location *models.LocationHistory, isPremium bool) error {
	// Check if a location history entry already exists for the user on the current day
	var existingLocation models.LocationHistory

	result := r.db.Where(
		`"Locationhistory"."UserID" = ? AND DATE_TRUNC('day', "Locationhistory"."Timestamp") = DATE_TRUNC('day', NOW())`,
		location.UserID,
	).First(&existingLocation)

	log.Print(result.RowsAffected)

	// If an entry exists and the user is not premium, return an error indicating that it cannot be created more than once a day
	if result.RowsAffected > 0 && !isPremium {
		return errors.New("location history already created for the user today")
	}

	// Create the location history entry
	result = r.db.Create(location)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *locationRepository) GetLocationHistoryByUserID(userID int) ([]models.LocationHistory, error) {
	var locationHistory []models.LocationHistory
	result := r.db.Where(`"Locationhistory"."UserID" = ?`, userID).Find(&locationHistory)
	if result.Error != nil {
		return nil, result.Error
	}
	return locationHistory, nil
}

func (r *locationRepository) GetNearbyLocations(userID int, maxDistance float64, page, pageSize int) ([]LocationWithDistance, error) {
	// Get the user's location
	var userLocation models.LocationHistory
	result := r.db.Where(`"Locationhistory"."UserID" = ?`, userID).Last(&userLocation)
	if result.Error != nil {
		return nil, result.Error
	}

	log.Print(userLocation.UserID, userLocation.Latitude, userLocation.Longitude)
	var nearbyLocations []LocationWithDistance

	// Get the profiles that have been shown to the user on the current day
	shownProfiles, err := r.getShownProfiles(userID)
	if err != nil {
		return nil, err
	}

	// Initialize shownProfiles as an empty slice if it's nil
	if shownProfiles == nil {
		shownProfiles = []int{}
	}

	// Conditionally include NOT IN clause
	var notInClause string
	var notInParams []interface{}
	if len(shownProfiles) > 0 {
		notInClause = `AND "Locationhistory"."UserID" NOT IN (?)`
		notInParams = append(notInParams, shownProfiles)
	}

	query := `
    SELECT
        "Locationhistory"."LocationID",
        "Locationhistory"."UserID",
        "Locationhistory"."Latitude",
        "Locationhistory"."Longitude",
        "Locationhistory"."Timestamp",
        (
            6371 * acos(
                cos(radians(?)) * cos(radians("Latitude")) * cos(radians("Longitude") - radians(?)) +
                sin(radians(?)) * sin(radians("Latitude"))
            )
        ) AS "distance"
    FROM "Locationhistory"
    WHERE "UserID" != ?
    AND (
        6371 * acos(
            cos(radians(?)) * cos(radians("Latitude")) * cos(radians("Longitude") - radians(?)) +
            sin(radians(?)) * sin(radians("Latitude"))
        )
    ) <= ?
`

	if len(shownProfiles) > 0 {
		query += notInClause
	}

	query += `
    ORDER BY "Timestamp" DESC
    LIMIT ? OFFSET ?;
`

	offset := (page - 1) * pageSize
	args := []interface{}{
		userLocation.Latitude, userLocation.Longitude, userLocation.Latitude,
		userID,
		userLocation.Latitude, userLocation.Longitude, userLocation.Latitude,
		maxDistance,
	}

	if len(shownProfiles) > 0 {
		args = append(args, notInParams...)
	}

	args = append(args, pageSize, offset)

	result = r.db.Raw(query, args...).Find(&nearbyLocations)

	// Log the query
	log.Printf("Executing query: %s\n", query)

	// Add logging
	log.Printf("Number of nearby locations: %d", len(nearbyLocations))
	log.Printf("Nearby locations: %+v", nearbyLocations)

	if result.Error != nil {
		return nil, result.Error
	}

	// Log the profiles that are being shown to the user
	err = r.logShownProfiles(userID, nearbyLocations)
	if err != nil {
		return nil, err
	}

	// Fetch profiles for the shown locations
	for i, loc := range nearbyLocations {
		profile, err := r.profileRepo.GetProfileByUserID(loc.UserID)
		if err == nil {
			nearbyLocations[i].Profile = profile
		}
	}

	return nearbyLocations, nil
}

func (r *locationRepository) getShownProfiles(userID int) ([]int, error) {
	var shownProfiles []int

	// Fetch the profiles that have been shown to the user on the current day
	err := r.db.Model(&models.ProfileView{}).
		Where(`"ViewerUserID" = ? AND DATE_TRUNC('day', "DateOnly") = DATE_TRUNC('day', NOW())`, userID).
		Pluck("ShownUserID", &shownProfiles).
		Error

	return shownProfiles, err
}

func (r *locationRepository) logShownProfiles(userID int, locations []LocationWithDistance) error {
	// Log the profiles that are being shown to the user on the current day
	var viewRecords []models.ProfileView
	currentTime := time.Now()

	for _, loc := range locations {
		viewRecord := models.ProfileView{
			ViewerUserID: userID,
			ShownUserID:  loc.UserID,
			Timestamp:    currentTime,
		}
		viewRecords = append(viewRecords, viewRecord)
	}

	// Print debug information
	log.Printf("Number of locations: %d", len(locations))
	log.Printf("Number of view records: %d", len(viewRecords))
	log.Printf("View records: %+v", viewRecords)

	// Exclude the "DateOnly" column from the insert operation
	err := r.db.Table("ProfileView").Omit("DateOnly").Create(&viewRecords).Error
	return err
}

// NewUserRepositoryWithGormDBAndRedis creates a new ProfileRepository with GormDB and Redis
func NewLocationRepositoryWithGormDBAndRedis(db helpers.DatabaseHandler, redis helpers.RedisHandler, profiles ProfileRepository) LocationRepository {
	return NewLocationRepository(db, redis, profiles)
}
