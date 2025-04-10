package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sangwan491/backend-assignments/employee-management/backend/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var collection *mongo.Collection
var validate *validator.Validate

func init() {
	// Initialize validator
	validate = validator.New()
}

// ConnectToMongoDB establishes a connection to MongoDB
// Returns an error if connection fails
func ConnectToMongoDB() error {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	// Get MongoDB connection details from environment variables
	connectionString := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DB_NAME")
	colName := os.Getenv("MONGODB_COLLECTION_NAME")

	if connectionString == "" || dbName == "" || colName == "" {
		return fmt.Errorf("missing required MongoDB environment variables")
	}

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return fmt.Errorf("MongoDB connection error: %w", err)
	}

	// Check the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("MongoDB ping error: %w", err)
	}

	collection = client.Database(dbName).Collection(colName)
	fmt.Println("MongoDB Connection success!")
	return nil
}

// formatValidationErrors converts validator errors into a user-friendly string.
func formatValidationErrors(errs validator.ValidationErrors) string {
	var errMsgs []string
	for _, err := range errs {
		// Provide more user-friendly messages based on the validation tag
		field := err.Field()
		tag := err.Tag()
		param := err.Param()

		var msg string
		switch tag {
		case "required":
			msg = fmt.Sprintf("Field '%s' is required", field)
		case "min":
			msg = fmt.Sprintf("Field '%s' must be at least %s", field, param)
		case "max":
			msg = fmt.Sprintf("Field '%s' must be at most %s", field, param)
		case "gt":
			msg = fmt.Sprintf("Field '%s' must be greater than %s", field, param)
		case "gte":
			msg = fmt.Sprintf("Field '%s' must be greater than or equal to %s", field, param)
		case "lt":
			msg = fmt.Sprintf("Field '%s' must be less than %s", field, param)
		case "lte":
			msg = fmt.Sprintf("Field '%s' must be less than or equal to %s", field, param)
		case "email":
			msg = fmt.Sprintf("Field '%s' must be a valid email address", field)
		// Add more cases for other common validation tags as needed
		default:
			msg = fmt.Sprintf("Field '%s' failed validation on the '%s' tag", field, tag)
		}
		errMsgs = append(errMsgs, msg)
	}
	return strings.Join(errMsgs, ", ")
}

// GetAllEmployees - HTTP handler to get all employees
func GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := getAllEmployees()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to retrieve employees: %v", err)})
		return
	}
	json.NewEncoder(w).Encode(employees)
}

// CreateEmployee - HTTP handler to create a new employee
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var employee models.Employee

	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Invalid request payload: %v", err)})
		return
	}

	if err := validate.Struct(employee); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": formatValidationErrors(validationErrors)})
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Validation error: %v", err)})
		return
	}

	if err := insertOneEmployee(employee); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to insert employee: %v", err)})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Employee created successfully"})
}

// UpdateEmployee - HTTP handler to update an employee
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	employeeID := params["id"]

	var employee models.Employee
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Invalid request payload: %v", err)})
		return
	}

	if err := validate.Struct(employee); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": formatValidationErrors(validationErrors)})
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Validation error: %v", err)})
		return
	}

	if err := updateOneEmployee(employeeID, employee); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to update employee: %v", err)})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Employee updated successfully"})
}

// DeleteEmployee - HTTP handler to delete an employee
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	employeeID := params["id"]

	if err := deleteOneEmployee(employeeID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to delete employee: %v", err)})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Employee deleted successfully"})
}

// insertOneEmployee inserts an employee into the database and returns an error if any.
func insertOneEmployee(employee models.Employee) error {
	result, err := collection.InsertOne(context.Background(), employee)
	if err != nil {
		return fmt.Errorf("error inserting employee: %w", err)
	}
	fmt.Println("Inserted 1 employee with id:", result.InsertedID)
	return nil
}

// updateOneEmployee updates an employee document in the database and returns an error if any.
func updateOneEmployee(employeeID string, employee models.Employee) error {
	id, err := bson.ObjectIDFromHex(employeeID)
	if err != nil {
		return fmt.Errorf("invalid employee ID format: %w", err)
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": employee}

	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("error updating employee: %w", err)
	}
	fmt.Println("Updated employee with id:", updateResult.UpsertedID)
	return nil
}

// deleteOneEmployee deletes an employee document from the database and returns an error if any.
func deleteOneEmployee(employeeID string) error {
	id, err := bson.ObjectIDFromHex(employeeID)
	if err != nil {
		return fmt.Errorf("invalid employee ID format: %w", err)
	}

	filter := bson.M{"_id": id}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("error deleting employee: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no employee found with ID: %s", employeeID)
	}
	fmt.Printf("Successfully deleted employee with ID: %s\n", employeeID)
	return nil
}

// getAllEmployees retrieves all employee documents from the database.
func getAllEmployees() ([]models.Employee, error) {
	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error finding employees: %w", err)
	}

	var employees []models.Employee
	for cur.Next(context.Background()) {
		var employee models.Employee
		if err := cur.Decode(&employee); err != nil {
			return nil, fmt.Errorf("error decoding employee: %w", err)
		}
		employees = append(employees, employee)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return employees, nil
}
