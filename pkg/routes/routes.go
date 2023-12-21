// routes/routes.go
package routes

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/metabbe3/knoxsdating/pkg/handlers"
	"github.com/metabbe3/knoxsdating/pkg/helpers"
	"github.com/metabbe3/knoxsdating/pkg/repository"
)

// InitializeRoutes initializes all routes for the application
func InitializeRoutes() *mux.Router {
	router := mux.NewRouter()

	// Connect to the database
	db, err := helpers.ConnectToDatabase()
	if err != nil {
		panic(err)
	}

	// Create repository instances
	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Update with your Redis server address
		// Add other Redis configuration options as needed
	})
	redisHelper := helpers.NewRedisHelper(redisClient)
	// Type assertion to *helpers.RedisHelper
	redisHelperInstance, ok := redisHelper.(*helpers.RedisHelper)
	if !ok {
		panic("Invalid type assertion for redisHelper")
	}

	// For Notification handlers
	notificationRepo := repository.NewNotificationRepositoryWithGormDBAndRedis(helpers.NewGormDBHandler(db), redisHelper)
	notificationHandlers := handlers.NewNotificationHandlers(notificationRepo, redisHelperInstance)

	// For User handlers
	userRepo := repository.NewUserRepositoryWithGormDBAndRedis(helpers.NewGormDBHandler(db), redisHelper)
	userHandlers := handlers.NewUserHandlers(userRepo, redisHelperInstance)

	// For Profile handlers
	profileRepo := repository.NewProfileRepositoryWithGormDBAndRedis(helpers.NewGormDBHandler(db), redisHelper)
	profileHandlers := handlers.NewProfileHandlers(profileRepo, redisHelperInstance)

	// For Location handlers
	locationRepo := repository.NewLocationRepositoryWithGormDBAndRedis(helpers.NewGormDBHandler(db), redisHelper, profileRepo)
	locationHandlers := handlers.NewLocationHandlers(locationRepo, userRepo, redisHelperInstance)

	// For SwipeHistory handlers
	swipeHistoryRepo := repository.NewSwipeHistoryRepositoryWithGormDBAndRedis(helpers.NewGormDBHandler(db), redisHelper)
	swipeHistoryHandlers := handlers.NewSwipeHistoryHandlers(swipeHistoryRepo, redisHelperInstance)
	// Add other handlers as needed

	router.HandleFunc("/notifications", notificationHandlers.CreateNotification).Methods("POST")
	router.HandleFunc("/notifications/{id:[0-9]+}", notificationHandlers.GetNotificationByID).Methods("GET")
	router.HandleFunc("/notifications/{id:[0-9]+}", notificationHandlers.UpdateNotification).Methods("PUT")
	router.HandleFunc("/notifications/{id:[0-9]+}", notificationHandlers.DeleteNotification).Methods("DELETE")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	}).Methods("GET")
	// Add other routes as needed

	router.HandleFunc("/users", userHandlers.RegisterUser).Methods("POST")
	router.HandleFunc("/users/login", userHandlers.Login).Methods("POST")
	router.HandleFunc("/users", userHandlers.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/updatePremium", userHandlers.UpdatePremium).Methods("PUT")

	// Profile routes
	router.HandleFunc("/profiles", profileHandlers.CreateProfile).Methods("POST")
	router.HandleFunc("/profiles", profileHandlers.GetProfile).Methods("GET")

	// Location routes
	router.HandleFunc("/locations", locationHandlers.CreateLocationHistory).Methods("POST")
	router.HandleFunc("/locations/nearby", locationHandlers.GetNearbyLocations).Methods("POST") // New route for getting nearby locations
	// Add other location routes as needed

	router.HandleFunc("/swipes", swipeHistoryHandlers.SaveSwipe).Methods("POST")
	router.HandleFunc("/swipes/matches", swipeHistoryHandlers.GetMatches).Methods("GET")
	router.HandleFunc("/swipes/redo", swipeHistoryHandlers.RedoSwipe).Methods("POST")

	return router
}
