package router

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sangwan491/backend-assignments/employee-management/backend/controllers"
)

// SetupRouter initializes all the routes for the application
func SetupRouter() http.Handler {
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api").Subrouter()

	// Employee routes
	api.HandleFunc("/employees", controllers.GetAllEmployees).Methods("GET")
	api.HandleFunc("/employees", controllers.CreateEmployee).Methods("POST")
	api.HandleFunc("/employees/{id}", controllers.UpdateEmployee).Methods("PUT")
	api.HandleFunc("/employees/{id}", controllers.DeleteEmployee).Methods("DELETE")

	return handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(router)
}
