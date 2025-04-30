package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sangwan491/backend-assignments/employee-management/backend/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var collection *mongo.Collection
var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ConnectToMongoDB() error {
	connectionString := "mongodb://mongodb:27017"
	dbName := "employee_db"
	colName := "employees"

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return fmt.Errorf("MongoDB connection error: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("MongoDB ping error: %w", err)
	}

	collection = client.Database(dbName).Collection(colName)
	fmt.Println("MongoDB Connection success!")

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"name", "text"},
			{"phone", "text"},
			{"email", "text"},
		},
	}

	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)

	if err != nil {
		return fmt.Errorf("Failed to create index: %w", err)
	}

	fmt.Println("Index created successfully!")

	return nil
}

func formatValidationErrors(errs validator.ValidationErrors) string {
	var errMsgs []string
	for _, err := range errs {
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
		default:
			msg = fmt.Sprintf("Field '%s' failed validation on the '%s' tag", field, tag)
		}
		errMsgs = append(errMsgs, msg)
	}
	return strings.Join(errMsgs, ", ")
}

// SearchEmployees - HTTP handler to search for employees with pagination
func SearchEmployees(w http.ResponseWriter, r *http.Request) {
	searchTerm := strings.TrimSpace(r.URL.Query().Get("term"))
	// searchType := r.URL.Query().Get("type")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Default to page 1 if invalid or missing
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	skip := int64((page - 1) * limit)

	pipeline := mongo.Pipeline{}

	if searchTerm != "" {
		matchStage := bson.D{{"$match", bson.M{"$text": bson.M{"$search": searchTerm}}}}

		pipeline = append(pipeline, matchStage)
	}

	facetStage := bson.D{{"$facet", bson.D{
		{"data", bson.A{
			bson.D{{"$skip", skip}},
			bson.D{{"$limit", limit}},
		}},
		{"totalCount", bson.A{
			bson.D{{"$count", "count"}},
		}},
	}}}

	pipeline = append(pipeline, facetStage)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cur, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Search query failed: %v", err)})
		return
	}
	defer cur.Close(ctx)

	var facetResults []struct {
		Data       []models.Employee `bson:"data"`
		TotalCount []struct {
			Count int64 `bson:"count"`
		} `bson:"totalCount"`
	}

	if err = cur.All(ctx, &facetResults); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to decode results: %v", err)})
		return
	}

	employees := []models.Employee{}
	if facetResults[0].Data != nil {
		employees = facetResults[0].Data
	}

	var totalCount int64 = 0
	if len(facetResults[0].TotalCount) > 0 {
		totalCount = facetResults[0].TotalCount[0].Count
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"employees": employees,
		"pages":     totalPages,
	})
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

	employeeID, err := insertOneEmployee(employee)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to insert employee: %v", err)})
		return
	}

	employee.ID = employeeID // Set the ID in the employee model

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Employee created successfully",
		"data":    employee,
	})
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
func insertOneEmployee(employee models.Employee) (bson.ObjectID, error) {
	result, err := collection.InsertOne(context.Background(), employee)
	if err != nil {
		return bson.NilObjectID, fmt.Errorf("error inserting employee: %w", err)
	}
	fmt.Println("Inserted 1 employee with id:", result.InsertedID)

	return result.InsertedID.(bson.ObjectID), nil // Return the inserted ID
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
