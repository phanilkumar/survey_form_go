package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open database
	db, err := sql.Open("sqlite3", "./survey_form.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Clear existing data
	_, err = db.Exec("DELETE FROM survey_responses")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("DELETE FROM surveys")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Creating sample surveys...")

	// Create sample surveys
	survey1Data := map[string]interface{}{
		"title":       "Customer Satisfaction Survey",
		"description": "Help us improve our services by providing your feedback on your recent experience.",
	}

	survey2Data := map[string]interface{}{
		"title":       "Employee Engagement Survey",
		"description": "We value your opinion! Please share your thoughts about workplace culture and satisfaction.",
	}

	survey3Data := map[string]interface{}{
		"title":       "Product Feedback Form",
		"description": "Tell us what you think about our latest product features and how we can make them better.",
	}

	// Insert surveys
	result1, err := db.Exec("INSERT INTO surveys (title, description) VALUES (?, ?)", survey1Data["title"], survey1Data["description"])
	if err != nil {
		log.Fatal(err)
	}
	survey1ID, _ := result1.LastInsertId()

	result2, err := db.Exec("INSERT INTO surveys (title, description) VALUES (?, ?)", survey2Data["title"], survey2Data["description"])
	if err != nil {
		log.Fatal(err)
	}
	survey2ID, _ := result2.LastInsertId()

	result3, err := db.Exec("INSERT INTO surveys (title, description) VALUES (?, ?)", survey3Data["title"], survey3Data["description"])
	if err != nil {
		log.Fatal(err)
	}
	survey3ID, _ := result3.LastInsertId()

	fmt.Println("Creating sample survey responses...")

	// Create sample responses for survey1
	response1Data := map[string]interface{}{
		"overall_satisfaction":      "5",
		"service_quality":           "4",
		"recommendation_likelihood": "5",
		"comments":                  "Great service, very satisfied!",
	}
	response1JSON, _ := json.Marshal(response1Data)

	response2Data := map[string]interface{}{
		"overall_satisfaction":      "3",
		"service_quality":           "4",
		"recommendation_likelihood": "3",
		"comments":                  "Service was okay, room for improvement.",
	}
	response2JSON, _ := json.Marshal(response2Data)

	// Create sample responses for survey2
	response3Data := map[string]interface{}{
		"workplace_culture":  "4",
		"job_satisfaction":   "5",
		"work_life_balance":  "4",
		"management_support": "5",
		"suggestions":        "More team building activities would be great!",
	}
	response3JSON, _ := json.Marshal(response3Data)

	// Create sample responses for survey3
	response4Data := map[string]interface{}{
		"product_rating":      "4",
		"feature_usefulness":  "5",
		"ease_of_use":         "4",
		"additional_features": "Mobile app would be helpful",
		"overall_impression":  "Very good product!",
	}
	response4JSON, _ := json.Marshal(response4Data)

	// Insert responses
	_, err = db.Exec("INSERT INTO survey_responses (survey_id, user_identifier, response_data) VALUES (?, ?, ?)", survey1ID, "user123", response1JSON)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO survey_responses (survey_id, user_identifier, response_data) VALUES (?, ?, ?)", survey1ID, "user456", response2JSON)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO survey_responses (survey_id, user_identifier, response_data) VALUES (?, ?, ?)", survey2ID, "employee001", response3JSON)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO survey_responses (survey_id, user_identifier, response_data) VALUES (?, ?, ?)", survey3ID, "customer789", response4JSON)
	if err != nil {
		log.Fatal(err)
	}

	// Count created data
	var surveyCount, responseCount int
	db.QueryRow("SELECT COUNT(*) FROM surveys").Scan(&surveyCount)
	db.QueryRow("SELECT COUNT(*) FROM survey_responses").Scan(&responseCount)

	fmt.Println("Sample data created successfully!")
	fmt.Printf("Created %d surveys and %d responses\n", surveyCount, responseCount)
}
