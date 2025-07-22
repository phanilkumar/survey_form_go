# 🐹 Survey Form API - Quick Reference

## **Base URL**
```
http://localhost:8080
```

## **API Endpoints**

### **📋 Survey Management**

#### **List All Surveys**
```http
GET /api/surveys
```

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "title": "Customer Satisfaction",
      "description": "Help us improve our services",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z",
      "responses_count": 2
    }
  ]
}
```

#### **Get Specific Survey**
```http
GET /api/surveys/{id}
```

#### **Create Survey**
```http
POST /api/surveys
Content-Type: application/json

{
  "survey": {
    "title": "New Survey",
    "description": "Survey description"
  }
}
```

**Validation:**
- Title: 3-255 characters
- Description: Required, max 1000 characters

### **📝 Survey Responses**

#### **List Survey Responses**
```http
GET /api/surveys/{survey_id}/survey_responses
```

#### **Get Specific Response**
```http
GET /api/surveys/{survey_id}/survey_responses/{id}
```

#### **Submit Response**
```http
POST /api/surveys/{survey_id}/survey_responses
Content-Type: application/json

{
  "survey_response": {
    "user_identifier": "john_doe",
    "response_data": {
      "rating": "5",
      "comment": "Great service!",
      "recommend": true
    }
  }
}
```

**Validation:**
- User Identifier: 3-100 characters
- Response Data: Required JSON object

#### **Update Response**
```http
PATCH /api/surveys/{survey_id}/survey_responses/{id}
Content-Type: application/json

{
  "survey_response": {
    "response_data": {
      "rating": "4",
      "comment": "Updated comment"
    }
  }
}
```

**Note:** Only editable within 24 hours of creation

### **👤 User Responses**

#### **Get User's Responses**
```http
GET /api/users/{user_identifier}/responses
```

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "survey": {
        "id": 1,
        "title": "Customer Satisfaction",
        "description": "Help us improve our services"
      },
      "user_identifier": "john_doe",
      "response_data": {
        "rating": "5",
        "comment": "Great service!"
      },
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z",
      "editable": true
    }
  ]
}
```

### **🔍 System Endpoints**

#### **API Information**
```http
GET /
```

#### **Health Check**
```http
GET /up
```

## **📊 Response Formats**

### **Success Response**
```json
{
  "status": "success",
  "message": "Operation completed successfully",
  "data": { ... }
}
```

### **Error Response**
```json
{
  "status": "error",
  "message": "Error description",
  "errors": [
    "Specific error message 1",
    "Specific error message 2"
  ]
}
```

## **🔢 HTTP Status Codes**

- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Invalid request data
- `404 Not Found` - Resource not found
- `422 Unprocessable Entity` - Validation errors
- `500 Internal Server Error` - Server error

## **🚨 Common Errors**

### **Validation Errors**
```json
{
  "status": "error",
  "message": "Failed to create survey",
  "errors": [
    "Title must be at least 3 characters long",
    "Description is required"
  ]
}
```

### **Not Found**
```json
{
  "status": "error",
  "message": "Survey not found"
}
```

### **Edit Time Expired**
```json
{
  "status": "error",
  "message": "Response cannot be edited after 24 hours"
}
```

## **🧪 Testing Examples**

### **Create and Test Survey**
```bash
# Create survey
curl -X POST http://localhost:8080/api/surveys \
  -H "Content-Type: application/json" \
  -d '{"survey":{"title":"Test Survey","description":"Test Description"}}'

# Submit response
curl -X POST http://localhost:8080/api/surveys/1/survey_responses \
  -H "Content-Type: application/json" \
  -d '{"survey_response":{"user_identifier":"testuser","response_data":{"rating":"5"}}}'

# Get user responses
curl http://localhost:8080/api/users/testuser/responses
```

## **📝 Notes**

- **Database**: SQLite (file: `survey_form.db`)
- **Port**: 8080 (configurable)
- **Response Data**: Flexible JSON structure
- **Editable Window**: 24 hours from creation
- **User Identifier**: Unique identifier for tracking responses 