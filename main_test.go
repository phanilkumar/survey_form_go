package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// TestSurvey represents a survey for testing
type TestSurvey struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	ResponsesCount int       `json:"responses_count"`
}

// TestSurveyResponse represents a survey response for testing
type TestSurveyResponse struct {
	ID             int             `json:"id"`
	SurveyID       int             `json:"survey_id"`
	UserIdentifier string          `json:"user_identifier"`
	ResponseData   json.RawMessage `json:"response_data"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Editable       bool            `json:"editable"`
}

// TestAPIResponse represents an API response for testing
type TestAPIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
}

var testDB *sql.DB

func setupTestDB() {
	var err error
	testDB, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	// Create tables
	createSurveysTable := `
	CREATE TABLE IF NOT EXISTS surveys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createResponsesTable := `
	CREATE TABLE IF NOT EXISTS survey_responses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		survey_id INTEGER NOT NULL,
		user_identifier TEXT NOT NULL,
		response_data TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (survey_id) REFERENCES surveys (id) ON DELETE CASCADE
	);`

	_, err = testDB.Exec(createSurveysTable)
	if err != nil {
		panic(err)
	}

	_, err = testDB.Exec(createResponsesTable)
	if err != nil {
		panic(err)
	}
}

func setupTestRouter() *gin.Engine {
	// Use test database
	db = testDB

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Root route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Survey Form API",
			"endpoints": gin.H{
				"surveys":        "/api/surveys",
				"responses":      "/api/surveys/{id}/responses",
				"user_responses": "/api/users/{user_identifier}/responses",
			},
		})
	})

	// API routes
	api := r.Group("/api")
	{
		api.GET("/surveys", getSurveys)
		api.POST("/surveys", createSurvey)
		api.GET("/surveys/:id", getSurvey)
		api.GET("/surveys/:id/responses", getSurveyResponses)
		api.POST("/surveys/:id/responses", createSurveyResponse)
		api.GET("/surveys/:id/responses/:response_id", getSurveyResponse)
		api.PATCH("/surveys/:id/responses/:response_id", updateSurveyResponse)
		api.GET("/users/:user_identifier/responses", getUserResponses)
	}

	return r
}

func TestGetSurveys(t *testing.T) {
	setupTestDB()
	defer testDB.Close()

	// Insert test data
	_, err := testDB.Exec("INSERT INTO surveys (title, description) VALUES (?, ?)", "Test Survey", "Test Description")
	assert.NoError(t, err)

	router := setupTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/surveys", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response TestAPIResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)

	// Check if data is present
	data, ok := response.Data.([]interface{})
	assert.True(t, ok)
	assert.Len(t, data, 1)
}

func TestGetSurvey(t *testing.T) {
	setupTestDB()
	defer testDB.Close()

	// Insert test survey
	result, err := testDB.Exec("INSERT INTO surveys (title, description) VALUES (?, ?)", "Test Survey", "Test Description")
	assert.NoError(t, err)
	surveyID, _ := result.LastInsertId()

	router := setupTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/surveys/%d", surveyID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response TestAPIResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
}

func TestGetSurveyNotFound(t *testing.T) {
	setupTestDB()
	defer testDB.Close()

	router := setupTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/surveys/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response TestAPIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, "Survey not found", response.Message)
}

func TestCreateSurvey(t *testing.T) {
	setupTestDB()
	defer testDB.Close()

	router := setupTestRouter()

	// Test valid survey creation
	surveyData := map[string]interface{}{
		"survey": map[string]string{
			"title":       "New Survey",
			"description": "New Description",
		},
	}

	jsonData, _ := json.Marshal(surveyData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/surveys", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response TestAPIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Survey created successfully", response.Message)
}

func TestCreateSurveyValidation(t *testing.T) {
	setupTestDB()
	defer testDB.Close()

	router := setupTestRouter()

	// Test invalid survey (short title)
	surveyData := map[string]interface{}{
		"survey": map[string]string{
			"title":       "A", // Too short
			"description": "New Description",
		},
	}

	jsonData, _ := json.Marshal(surveyData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/surveys", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var response TestAPIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)
	assert.Contains(t, response.Errors[0], "Title must be at least 3 characters")
}

func TestCreateSurveyResponse(t *testing.T) {
	setupTestDB()
	defer testDB.Close()

	// Create a survey first
	result, err := testDB.Exec("INSERT INTO surveys (title, description) VALUES (?, ?)", "Test Survey", "Test Description")
	assert.NoError(t, err)
	surveyID, _ := result.LastInsertId()

	router := setupTestRouter()

	// Test valid response creation
	responseData := map[string]interface{}{
		"survey_response": map[string]interface{}{
			"user_identifier": "testuser",
			"response_data":   json.RawMessage(`{"rating": "5", "comment": "Great!"}`),
		},
	}

	jsonData, _ := json.Marshal(responseData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/surveys/%d/responses", surveyID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response TestAPIResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Survey response submitted successfully", response.Message)
}

func TestCreateSurveyResponseValidation(t *testing.T) {
	setupTestDB()
	defer testDB.Close()

	// Create a survey first
	result, err := testDB.Exec("INSERT INTO surveys (title, description) VALUES (?, ?)", "Test Survey", "Test Description")
	assert.NoError(t, err)
	surveyID, _ := result.LastInsertId()

	router := setupTestRouter()

	// Test invalid response (short user identifier)
	responseData := map[string]interface{}{
		"survey_response": map[string]interface{}{
			"user_identifier": "ab", // Too short
			"response_data":   json.RawMessage(`{"rating": "5"}`),
		},
	}

	jsonData, _ := json.Marshal(responseData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/surveys/%d/responses", surveyID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var response TestAPIResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)
	assert.Contains(t, response.Errors[0], "User identifier must be at least 3 characters")
}

func TestGetSurveyResponses(t *testing.T) {
	setupTestDB()
	defer testDB.Close()

	// Create a survey and response
	result, err := testDB.Exec("INSERT INTO surveys (title, description) VALUES (?, ?)", "Test Survey", "Test Description")
	assert.NoError(t, err)
	surveyID, _ := result.LastInsertId()

	responseData := json.RawMessage(`{"rating": "5"}`)
	_, err = testDB.Exec("INSERT INTO survey_responses (survey_id, user_identifier, response_data) VALUES (?, ?, ?)", surveyID, "testuser", responseData)
	assert.NoError(t, err)

	router := setupTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/surveys/%d/responses", surveyID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response TestAPIResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)

	// Check if data is present
	data, ok := response.Data.([]interface{})
	assert.True(t, ok)
	assert.Len(t, data, 1)
}

func TestGetUserResponses(t *testing.T) {
	setupTestDB()
	defer testDB.Close()

	// Create a survey and response
	result, err := testDB.Exec("INSERT INTO surveys (title, description) VALUES (?, ?)", "Test Survey", "Test Description")
	assert.NoError(t, err)
	surveyID, _ := result.LastInsertId()

	responseData := json.RawMessage(`{"rating": "5"}`)
	_, err = testDB.Exec("INSERT INTO survey_responses (survey_id, user_identifier, response_data) VALUES (?, ?, ?)", surveyID, "testuser", responseData)
	assert.NoError(t, err)

	router := setupTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/users/testuser/responses", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response TestAPIResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)

	// Check if data is present
	data, ok := response.Data.([]interface{})
	assert.True(t, ok)
	assert.Len(t, data, 1)
}

func TestUpdateSurveyResponse(t *testing.T) {
	setupTestDB()
	defer testDB.Close()

	// Create a survey and response
	result, err := testDB.Exec("INSERT INTO surveys (title, description) VALUES (?, ?)", "Test Survey", "Test Description")
	assert.NoError(t, err)
	surveyID, _ := result.LastInsertId()

	responseData := json.RawMessage(`{"rating": "5"}`)
	result2, err := testDB.Exec("INSERT INTO survey_responses (survey_id, user_identifier, response_data) VALUES (?, ?, ?)", surveyID, "testuser", responseData)
	assert.NoError(t, err)
	responseID, _ := result2.LastInsertId()

	router := setupTestRouter()

	// Test valid response update
	updateData := map[string]interface{}{
		"survey_response": map[string]interface{}{
			"response_data": json.RawMessage(`{"rating": "4", "comment": "Updated!"}`),
		},
	}

	jsonData, _ := json.Marshal(updateData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/api/surveys/%d/responses/%d", surveyID, responseID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response TestAPIResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Survey response updated successfully", response.Message)
}

func TestRootEndpoint(t *testing.T) {
	setupTestDB()
	defer testDB.Close()

	router := setupTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Survey Form API", response["message"])
}
