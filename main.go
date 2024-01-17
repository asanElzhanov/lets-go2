package main

import (
	"encoding/json"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
type User struct {
    gorm.Model
    Name string
    Email string
}


func main() {
	connStr := "user=username dbname=mydb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
	log.Fatal(err)
		}
	err = db.Ping()
	if err != nil {
	log.Fatal(err)
		}
	fmt.Println("Successfully connected to the database")
	defer db.Close()

	// ###############################################

	http.HandleFunc("/receiveData", receiveData)
	http.HandleFunc("/handleGetRequest", handleGetRequest)
	http.HandleFunc("/htmlPage", handleHTMLPage)

	port := 8080
	serverAddr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server is listening on http://localhost%s\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func receiveData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
        http.Error(w, "ERROR", http.StatusMethodNotAllowed)
        return
    }

    var requestData map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        http.Error(w, "ERROR", http.StatusBadRequest)
        return
    }

    person, ok := requestData["person"].(map[string]interface{})
    if !ok {
        handleErrorResponse(w, "Неверный формат JSON")
        return
    }
    response := Response{
        Status:  "success",
        Message: "SUCCESS",
    }
    sendJSONResponse(w, http.StatusOK, response)
}
func handleErrorResponse(w http.ResponseWriter, errorMessage string) {
    response := Response{
        Status:  "400",
        Message: errorMessage,
    }
    sendJSONResponse(w, http.StatusBadRequest, response)
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}
func handleGetRequest(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "ERROR", http.StatusMethodNotAllowed)
        return
    }

    personID := r.URL.Query().Get("id")
    name := r.URL.Query().Get("name")

    response := Response{
        Status:  "success",
        Message: fmt.Sprintf("Person ID: %s, Name: %s", personID, name),
    }
    sendJSONResponse(w, http.StatusOK, response)
}
func handleHTMLPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
func createUser(db *gorm.DB, user *User) error {
    result := db.Create(user)
    if result.Error != nil {
        return result.Error
    }
    return nil
}
func getUserByID(db *gorm.DB, userID uint) (*User, error) {
    var user User
    result := db.First(&user, userID)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}
func updateUserName(db *gorm.DB, userID uint, newName string) error {
    result := db.Model(&User{}).Where("id = ?", userID).Update("name", newName)
    if result.Error != nil {
        return result.Error
    }
    return nil
}
func deleteUserByID(db *gorm.DB, userID uint) error {
    result := db.Delete(&User{}, userID)
    if result.Error != nil {
        return result.Error
    }
    return nil
}
func getAllUsers(db *gorm.DB) ([]User, error) {
    var users []User
    result := db.Find(&users)
    if result.Error != nil {
        return nil, result.Error
    }
    return users, nil
}
