package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// Survey represents a survey in the database
type Survey struct {
	ID             int       `json:"id" db:"id"`
	Title          string    `json:"title" db:"title"`
	Description    string    `json:"description" db:"description"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	ResponsesCount int       `json:"responses_count"`
}

// SurveyResponse represents a survey response in the database
type SurveyResponse struct {
	ID             int             `json:"id" db:"id"`
	SurveyID       int             `json:"survey_id" db:"survey_id"`
	UserIdentifier string          `json:"user_identifier" db:"user_identifier"`
	ResponseData   json.RawMessage `json:"response_data" db:"response_data"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
	Editable       bool            `json:"editable"`
}

// UserResponse represents a response with survey information
type UserResponse struct {
	ID             int             `json:"id"`
	Survey         Survey          `json:"survey"`
	UserIdentifier string          `json:"user_identifier"`
	ResponseData   json.RawMessage `json:"response_data"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Editable       bool            `json:"editable"`
}

// CreateSurveyRequest represents the request body for creating a survey
type CreateSurveyRequest struct {
	Survey struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
	} `json:"survey" binding:"required"`
}

// CreateResponseRequest represents the request body for creating a response
type CreateResponseRequest struct {
	SurveyResponse struct {
		UserIdentifier string          `json:"user_identifier" binding:"required"`
		ResponseData   json.RawMessage `json:"response_data" binding:"required"`
	} `json:"survey_response" binding:"required"`
}

// UpdateResponseRequest represents the request body for updating a response
type UpdateResponseRequest struct {
	SurveyResponse struct {
		ResponseData json.RawMessage `json:"response_data"`
	} `json:"survey_response"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
}

// Database connection
var db *sql.DB

func main() {
	// Initialize database
	initDatabase()
	defer db.Close()

	// Create Gin router
	r := gin.Default()

	// API routes
	api := r.Group("/api")
	{
		// Survey routes
		api.GET("/surveys", getSurveys)
		api.POST("/surveys", createSurvey)
		api.GET("/surveys/:id", getSurvey)

		// Survey response routes
		api.GET("/surveys/:id/responses", getSurveyResponses)
		api.POST("/surveys/:id/responses", createSurveyResponse)
		api.GET("/surveys/:id/responses/:response_id", getSurveyResponse)
		api.PATCH("/surveys/:id/responses/:response_id", updateSurveyResponse)

		// User response routes
		api.GET("/users/:user_identifier/responses", getUserResponses)
	}

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

	// Health check
	r.GET("/up", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Run the server
	fmt.Println("Server running on http://localhost:8081")
	r.Run(":8081")
}

// initDatabase initializes the SQLite database and creates tables
func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./survey_form.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create surveys table
	createSurveysTable := `
	CREATE TABLE IF NOT EXISTS surveys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Create survey_responses table
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

	_, err = db.Exec(createSurveysTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createResponsesTable)
	if err != nil {
		log.Fatal(err)
	}
}

// getSurveys returns all surveys
func getSurveys(c *gin.Context) {
	rows, err := db.Query(`
		SELECT s.id, s.title, s.description, s.created_at, s.updated_at,
		       COUNT(sr.id) as responses_count
		FROM surveys s
		LEFT JOIN survey_responses sr ON s.id = sr.survey_id
		GROUP BY s.id
		ORDER BY s.created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch surveys",
			Errors:  []string{err.Error()},
		})
		return
	}
	defer rows.Close()

	var surveys []Survey
	for rows.Next() {
		var survey Survey
		err := rows.Scan(&survey.ID, &survey.Title, &survey.Description, &survey.CreatedAt, &survey.UpdatedAt, &survey.ResponsesCount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Status:  "error",
				Message: "Failed to scan survey data",
				Errors:  []string{err.Error()},
			})
			return
		}
		surveys = append(surveys, survey)
	}

	c.JSON(http.StatusOK, APIResponse{
		Status: "success",
		Data:   surveys,
	})
}

// getSurvey returns a specific survey
func getSurvey(c *gin.Context) {
	id := c.Param("id")
	surveyID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid survey ID",
			Errors:  []string{err.Error()},
		})
		return
	}

	var survey Survey
	err = db.QueryRow(`
		SELECT s.id, s.title, s.description, s.created_at, s.updated_at,
		       COUNT(sr.id) as responses_count
		FROM surveys s
		LEFT JOIN survey_responses sr ON s.id = sr.survey_id
		WHERE s.id = ?
		GROUP BY s.id
	`, surveyID).Scan(&survey.ID, &survey.Title, &survey.Description, &survey.CreatedAt, &survey.UpdatedAt, &survey.ResponsesCount)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, APIResponse{
				Status:  "error",
				Message: "Survey not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch survey",
			Errors:  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status: "success",
		Data:   survey,
	})
}

// createSurvey creates a new survey
func createSurvey(c *gin.Context) {
	var req CreateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid request data",
			Errors:  []string{err.Error()},
		})
		return
	}

	// Validation
	var errors []string
	if len(req.Survey.Title) < 3 {
		errors = append(errors, "Title must be at least 3 characters long")
	}
	if len(req.Survey.Title) > 255 {
		errors = append(errors, "Title must be less than 255 characters")
	}
	if len(req.Survey.Description) > 1000 {
		errors = append(errors, "Description must be less than 1000 characters")
	}

	if len(errors) > 0 {
		c.JSON(http.StatusUnprocessableEntity, APIResponse{
			Status:  "error",
			Message: "Failed to create survey",
			Errors:  errors,
		})
		return
	}

	result, err := db.Exec(`
		INSERT INTO surveys (title, description, created_at, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, req.Survey.Title, req.Survey.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to create survey",
			Errors:  []string{err.Error()},
		})
		return
	}

	id, _ := result.LastInsertId()
	var survey Survey
	err = db.QueryRow(`
		SELECT id, title, description, created_at, updated_at, 0 as responses_count
		FROM surveys WHERE id = ?
	`, id).Scan(&survey.ID, &survey.Title, &survey.Description, &survey.CreatedAt, &survey.UpdatedAt, &survey.ResponsesCount)

	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch created survey",
			Errors:  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Status:  "success",
		Message: "Survey created successfully",
		Data:    survey,
	})
}

// getSurveyResponses returns all responses for a survey
func getSurveyResponses(c *gin.Context) {
	surveyID := c.Param("id")
	id, err := strconv.Atoi(surveyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid survey ID",
			Errors:  []string{err.Error()},
		})
		return
	}

	// Check if survey exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM surveys WHERE id = ?)", id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "Survey not found",
		})
		return
	}

	rows, err := db.Query(`
		SELECT id, survey_id, user_identifier, response_data, created_at, updated_at
		FROM survey_responses
		WHERE survey_id = ?
		ORDER BY updated_at DESC
	`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch responses",
			Errors:  []string{err.Error()},
		})
		return
	}
	defer rows.Close()

	var responses []SurveyResponse
	for rows.Next() {
		var response SurveyResponse
		err := rows.Scan(&response.ID, &response.SurveyID, &response.UserIdentifier, &response.ResponseData, &response.CreatedAt, &response.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Status:  "error",
				Message: "Failed to scan response data",
				Errors:  []string{err.Error()},
			})
			return
		}
		response.Editable = time.Since(response.CreatedAt) < 24*time.Hour
		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, APIResponse{
		Status: "success",
		Data:   responses,
	})
}

// getSurveyResponse returns a specific survey response
func getSurveyResponse(c *gin.Context) {
	surveyID := c.Param("id")
	responseID := c.Param("response_id")

	sID, err := strconv.Atoi(surveyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid survey ID",
			Errors:  []string{err.Error()},
		})
		return
	}

	rID, err := strconv.Atoi(responseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid response ID",
			Errors:  []string{err.Error()},
		})
		return
	}

	var response SurveyResponse
	err = db.QueryRow(`
		SELECT id, survey_id, user_identifier, response_data, created_at, updated_at
		FROM survey_responses
		WHERE id = ? AND survey_id = ?
	`, rID, sID).Scan(&response.ID, &response.SurveyID, &response.UserIdentifier, &response.ResponseData, &response.CreatedAt, &response.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, APIResponse{
				Status:  "error",
				Message: "Survey response not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch response",
			Errors:  []string{err.Error()},
		})
		return
	}

	response.Editable = time.Since(response.CreatedAt) < 24*time.Hour

	c.JSON(http.StatusOK, APIResponse{
		Status: "success",
		Data:   response,
	})
}

// createSurveyResponse creates a new survey response
func createSurveyResponse(c *gin.Context) {
	surveyID := c.Param("id")
	sID, err := strconv.Atoi(surveyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid survey ID",
			Errors:  []string{err.Error()},
		})
		return
	}

	// Check if survey exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM surveys WHERE id = ?)", sID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "Survey not found",
		})
		return
	}

	var req CreateResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid request data",
			Errors:  []string{err.Error()},
		})
		return
	}

	// Validation
	var errors []string
	if len(req.SurveyResponse.UserIdentifier) < 3 {
		errors = append(errors, "User identifier must be at least 3 characters long")
	}
	if len(req.SurveyResponse.UserIdentifier) > 100 {
		errors = append(errors, "User identifier must be less than 100 characters")
	}

	if len(errors) > 0 {
		c.JSON(http.StatusUnprocessableEntity, APIResponse{
			Status:  "error",
			Message: "Failed to submit survey response",
			Errors:  errors,
		})
		return
	}

	result, err := db.Exec(`
		INSERT INTO survey_responses (survey_id, user_identifier, response_data, created_at, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, sID, req.SurveyResponse.UserIdentifier, req.SurveyResponse.ResponseData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to submit survey response",
			Errors:  []string{err.Error()},
		})
		return
	}

	id, _ := result.LastInsertId()
	var response SurveyResponse
	err = db.QueryRow(`
		SELECT id, survey_id, user_identifier, response_data, created_at, updated_at
		FROM survey_responses WHERE id = ?
	`, id).Scan(&response.ID, &response.SurveyID, &response.UserIdentifier, &response.ResponseData, &response.CreatedAt, &response.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch created response",
			Errors:  []string{err.Error()},
		})
		return
	}

	response.Editable = true

	c.JSON(http.StatusCreated, APIResponse{
		Status:  "success",
		Message: "Survey response submitted successfully",
		Data:    response,
	})
}

// updateSurveyResponse updates a survey response
func updateSurveyResponse(c *gin.Context) {
	surveyID := c.Param("id")
	responseID := c.Param("response_id")

	sID, err := strconv.Atoi(surveyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid survey ID",
			Errors:  []string{err.Error()},
		})
		return
	}

	rID, err := strconv.Atoi(responseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid response ID",
			Errors:  []string{err.Error()},
		})
		return
	}

	// Check if response exists and is editable
	var response SurveyResponse
	err = db.QueryRow(`
		SELECT id, survey_id, user_identifier, response_data, created_at, updated_at
		FROM survey_responses
		WHERE id = ? AND survey_id = ?
	`, rID, sID).Scan(&response.ID, &response.SurveyID, &response.UserIdentifier, &response.ResponseData, &response.CreatedAt, &response.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, APIResponse{
				Status:  "error",
				Message: "Survey response not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch response",
			Errors:  []string{err.Error()},
		})
		return
	}

	// Check if response is editable (within 24 hours)
	if time.Since(response.CreatedAt) >= 24*time.Hour {
		c.JSON(http.StatusUnprocessableEntity, APIResponse{
			Status:  "error",
			Message: "Response cannot be edited after 24 hours",
		})
		return
	}

	var req UpdateResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid request data",
			Errors:  []string{err.Error()},
		})
		return
	}

	// Update response data
	_, err = db.Exec(`
		UPDATE survey_responses 
		SET response_data = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND survey_id = ?
	`, req.SurveyResponse.ResponseData, rID, sID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to update survey response",
			Errors:  []string{err.Error()},
		})
		return
	}

	// Fetch updated response
	err = db.QueryRow(`
		SELECT id, survey_id, user_identifier, response_data, created_at, updated_at
		FROM survey_responses WHERE id = ?
	`, rID).Scan(&response.ID, &response.SurveyID, &response.UserIdentifier, &response.ResponseData, &response.CreatedAt, &response.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch updated response",
			Errors:  []string{err.Error()},
		})
		return
	}

	response.Editable = time.Since(response.CreatedAt) < 24*time.Hour

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Survey response updated successfully",
		Data:    response,
	})
}

// getUserResponses returns all responses for a specific user
func getUserResponses(c *gin.Context) {
	userIdentifier := c.Param("user_identifier")

	rows, err := db.Query(`
		SELECT sr.id, sr.survey_id, sr.user_identifier, sr.response_data, sr.created_at, sr.updated_at,
		       s.id, s.title, s.description
		FROM survey_responses sr
		JOIN surveys s ON sr.survey_id = s.id
		WHERE sr.user_identifier = ?
		ORDER BY sr.updated_at DESC
	`, userIdentifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch user responses",
			Errors:  []string{err.Error()},
		})
		return
	}
	defer rows.Close()

	var responses []UserResponse
	for rows.Next() {
		var response UserResponse
		var survey Survey
		err := rows.Scan(&response.ID, &response.Survey.ID, &response.UserIdentifier, &response.ResponseData, &response.CreatedAt, &response.UpdatedAt, &survey.ID, &survey.Title, &survey.Description)
		if err != nil {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Status:  "error",
				Message: "Failed to scan user response data",
				Errors:  []string{err.Error()},
			})
			return
		}
		response.Survey = survey
		response.Editable = time.Since(response.CreatedAt) < 24*time.Hour
		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, APIResponse{
		Status: "success",
		Data:   responses,
	})
}
