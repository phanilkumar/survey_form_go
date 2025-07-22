# 🐹 Survey Form API - Go (Golang) Implementation

A complete **RESTful API** for managing surveys and responses, built with **Go (Golang)** and the **Gin framework**. Perfect for beginners learning Go web development!

## 🎯 **Features**

### ✅ **Basic Requirements**
- [x] **Survey Management**: Create and retrieve surveys
- [x] **Response Submission**: Submit responses to surveys
- [x] **Response Updates**: Edit responses within 24 hours
- [x] **User Response History**: View all responses by user identifier
- [x] **RESTful API**: Clean, standard REST endpoints

### ✅ **Next-Level Requirements**
- [x] **Data Validation**: Comprehensive input validation
- [x] **Error Handling**: Detailed error messages and status codes
- [x] **Response Count**: Track number of responses per survey
- [x] **Editable Responses**: 24-hour editing window
- [x] **JSON Response Data**: Flexible response data storage
- [x] **SQLite Database**: Lightweight, file-based database

## 🚀 **Quick Start**

### **Prerequisites**
- **Go 1.21+** installed on your system
- Basic knowledge of Go syntax (helpful but not required)

### **Installation & Setup**

1. **Clone or download the project**
   ```bash
   git clone <repository-url>
   cd survey_form_go
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Run the application**
   ```bash
   go run main.go
   ```

4. **Populate with sample data** (optional)
   ```bash
   go run seed.go
   ```

5. **Test the API**
   ```bash
   curl http://localhost:8081/
   ```

## 📚 **API Endpoints**

### **Root & Health Check**
- `GET /` - API information and available endpoints
- `GET /up` - Health check endpoint

### **Survey Management**
- `GET /api/surveys` - List all surveys
- `GET /api/surveys/:id` - Get specific survey details
- `POST /api/surveys` - Create a new survey

### **Survey Responses**
- `GET /api/surveys/:id/responses` - List all responses for a survey
- `GET /api/surveys/:id/responses/:response_id` - Get specific response
- `POST /api/surveys/:id/responses` - Submit a new response
- `PATCH /api/surveys/:id/responses/:response_id` - Update an existing response

### **User Responses**
- `GET /api/users/:user_identifier/responses` - Get all responses by a user

## 🔧 **Usage Examples**

### **Create a Survey**
```bash
curl -X POST http://localhost:8080/api/surveys \
  -H "Content-Type: application/json" \
  -d '{
    "survey": {
      "title": "Customer Feedback",
      "description": "Help us improve our services"
    }
  }'
```

### **Submit a Response**
```bash
curl -X POST http://localhost:8081/api/surveys/1/responses \
  -H "Content-Type: application/json" \
  -d '{
    "survey_response": {
      "user_identifier": "john_doe",
      "response_data": {
        "rating": "5",
        "comment": "Excellent service!",
        "recommend": true
      }
    }
  }'
```

### **Get User Responses**
```bash
curl http://localhost:8081/api/users/john_doe/responses
```

### **Update a Response**
```bash
curl -X PATCH http://localhost:8081/api/surveys/1/responses/1 \
  -H "Content-Type: application/json" \
  -d '{
    "survey_response": {
      "response_data": {
        "rating": "4",
        "comment": "Updated comment",
        "recommend": true
      }
    }
  }'
```

## 🧪 **Testing**

Run the comprehensive test suite:

```bash
go test -v
```

### **Test Coverage**
- ✅ Survey creation and retrieval
- ✅ Response submission and updates
- ✅ Input validation
- ✅ Error handling
- ✅ User response history
- ✅ API endpoints

## 📊 **Database Schema**

### **Surveys Table**
```sql
CREATE TABLE surveys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### **Survey Responses Table**
```sql
CREATE TABLE survey_responses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    survey_id INTEGER NOT NULL,
    user_identifier TEXT NOT NULL,
    response_data TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (survey_id) REFERENCES surveys (id) ON DELETE CASCADE
);
```

## 🔍 **Response Format**

### **Success Response**
```json
{
  "status": "success",
  "message": "Survey created successfully",
  "data": {
    "id": 1,
    "title": "Customer Feedback",
    "description": "Help us improve our services",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "responses_count": 0
  }
}
```

### **Error Response**
```json
{
  "status": "error",
  "message": "Failed to create survey",
  "errors": [
    "Title must be at least 3 characters long"
  ]
}
```

## 🛠 **Project Structure**

```
survey_form_go/
├── main.go              # Main application file
├── main_test.go         # Comprehensive test suite
├── seed_data.go         # Sample data population
├── go.mod              # Go module dependencies
├── go.sum              # Dependency checksums
├── README.md           # This file
└── survey_form.db      # SQLite database (created automatically)
```

## 🎓 **Learning Go with This Project**

### **Key Go Concepts Covered**
- **Structs**: Data models for surveys and responses
- **Interfaces**: HTTP handlers and database operations
- **Error Handling**: Comprehensive error management
- **JSON Marshaling**: Request/response serialization
- **Database Operations**: SQL queries with Go's database/sql
- **HTTP Routing**: RESTful API with Gin framework
- **Testing**: Unit tests with testify framework

### **Why This is Great for Beginners**
1. **Simple Structure**: Easy to understand and modify
2. **Real-world API**: Practical RESTful implementation
3. **Clear Comments**: Well-documented code
4. **Comprehensive Testing**: Learn testing practices
5. **Step-by-step Setup**: Detailed installation instructions
6. **Modern Go**: Uses current Go best practices

## 🔧 **Configuration**

### **Database**
- **Type**: SQLite (file-based, no setup required)
- **File**: `survey_form.db` (created automatically)
- **Location**: Project root directory

### **Server**
- **Port**: 8081 (configurable in main.go)
- **Host**: localhost
- **URL**: http://localhost:8081

## 🚨 **Validation Rules**

### **Survey Creation**
- Title: 3-255 characters
- Description: Required, max 1000 characters

### **Response Submission**
- User Identifier: 3-100 characters
- Response Data: Required JSON object
- Survey must exist

### **Response Updates**
- Only editable within 24 hours of creation
- Response must exist and belong to specified survey

## 🐛 **Troubleshooting**

### **Common Issues**

1. **"command not found: go"**
   - Install Go from https://golang.org/dl/
   - Add Go to your PATH

2. **"module not found" errors**
   - Run `go mod tidy` to download dependencies
   - Check your internet connection

3. **"port already in use"**
   - Change the port in main.go (line with `r.Run(":8080")`)
   - Or kill the process using port 8080

4. **Database errors**
   - Delete `survey_form.db` and restart the application
   - Check file permissions in the project directory

### **Getting Help**
- Check the test files for usage examples
- Review the API documentation above
- Look at the sample data in `seed_data.go`

## 🎉 **Success!**

You now have a complete, production-ready survey form API built with Go! This implementation includes:

- ✅ **All basic requirements** implemented
- ✅ **All next-level requirements** implemented  
- ✅ **Comprehensive testing** (13 test cases)
- ✅ **Sample data** for immediate testing
- ✅ **Full documentation** and examples
- ✅ **Beginner-friendly** code structure

## 🔗 **Related Projects**

This is part of a series of survey form implementations:
- **Ruby on Rails**: https://github.com/phanilkumar/survey_form.git
- **Python Flask**: https://github.com/phanilkumar/survey_form_python.git
- **Go Gin**: This project

Each implementation demonstrates the same functionality in different languages and frameworks, perfect for learning and comparison!

---

**Happy coding! 🐹✨** 