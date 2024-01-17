package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
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

    // Доступ к другим полям в структуре 'person' по мере необходимости
    // Например:
    // personID := person["id"].(float64)
    // name := person["name"].(string)

    // Обработка данных и ответ успешным статусом
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

    // Доступ к параметрам запроса
    personID := r.URL.Query().Get("id")
    name := r.URL.Query().Get("name")

    // Обработка параметров и ответ успешным статусом
    response := Response{
        Status:  "success",
        Message: fmt.Sprintf("Person ID: %s, Name: %s", personID, name),
    }
    sendJSONResponse(w, http.StatusOK, response)
}

